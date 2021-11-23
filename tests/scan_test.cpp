#include <sstream>

#include <benchmark/benchmark.h>

#include <hs/hs.h>

#include "scan_test.h"

#if HAVE_RE2
#include <re2/re2.h>
using namespace re2;
#endif

#if HAVE_PCRE2
#define PCRE2_CODE_UNIT_WIDTH 8
#define PCRE2_STATIC 1
#include <pcre2.h>
#endif

enum benchCase
{
    Easy0,
    Easy0i,
    Easy1,
    Medium,
    Hard,
    Hard1,
};

std::map<benchCase, std::string> benchData{
    {Easy0, "ABCDEFGHIJKLMNOPQRSTUVWXYZ$"},
    {Easy0i, "(?i)ABCDEFGHIJklmnopqrstuvwxyz$"},
    {Easy1, "A[AB]B[BC]C[CD]D[DE]E[EF]F[FG]G[GH]H[HI]I[IJ]J$"},
    {Medium, "[XYZ]ABCDEFGHIJKLMNOPQRSTUVWXYZ$"},
    {Hard, "[ -~]*ABCDEFGHIJKLMNOPQRSTUVWXYZ$"},
    {Hard1, "ABCD|CDEF|EFGH|GHIJ|IJKL|KLMN|MNOP|OPQR|QRST|STUV|UVWX|WXYZ"},
};

std::vector<std::vector<int64_t> > args = {
    {
        Easy0,
        Easy0i,
        Easy1,
        Medium,
        Hard,
        Hard1,
    },
    {
        16,
        32,
        1 << 10,
        32 << 10,
        1 << 20,
        32 << 20,
    },
};

std::vector<char> make_text(int n)
{
    std::vector<char> text(n);

    uint32_t x = ~0;

    for (int i = 0; i < n; i++)
    {
        x += x;
        x ^= 1;
        if (int(x) < 0)
        {
            x ^= 0x88888eef;
        }
        if (x % 31 == 0)
        {
            text[i] = '\n';
        }
        else
        {
            text[i] = x % (0x7E + 1 - 0x20) + 0x20;
        }
    }
    return text;
}

int on_match_event(unsigned int id, unsigned long long from, unsigned long long to, unsigned int flags, void *context)
{
    return 0;
}

static void BM_BlockScan(benchmark::State &state)
{
    hs_database_t *db = nullptr;
    hs_compile_error_t *compile_err = nullptr;
    hs_scratch_t *s = nullptr;

    auto expr = benchData[benchCase(state.range(0))];

    if (hs_compile(expr.c_str(), HS_FLAG_MULTILINE, HS_MODE_BLOCK, nullptr, &db, &compile_err) != HS_SUCCESS)
    {
        state.SkipWithError("compile failed");
    }

    if (hs_alloc_scratch(db, &s) != HS_SUCCESS)
    {
        state.SkipWithError("alloc scratch");
    }

    auto text = make_text(state.range(1));

    for (auto _ : state)
    {
        if (hs_scan(db, text.data(), text.size(), 0, s, on_match_event, nullptr) != HS_SUCCESS)
        {
            state.SkipWithError("scan failed");
        }
    }

    state.SetBytesProcessed(int64_t(state.iterations()) * int64_t(state.range(1)));

    hs_free_scratch(s);
}

BENCHMARK(BM_BlockScan)->ArgsProduct(args)->ArgNames({"regex", "size"});

const size_t page_size = 4096;

static void BM_StreamScan(benchmark::State &state)
{
    hs_database_t *db = nullptr;
    hs_compile_error_t *compile_err = nullptr;
    hs_scratch_t *s = nullptr;

    auto expr = benchData[benchCase(state.range(0))];

    if (hs_compile(expr.c_str(), HS_FLAG_MULTILINE, HS_MODE_STREAM, nullptr, &db, &compile_err) != HS_SUCCESS)
    {
        state.SkipWithError("compile failed");
    }

    if (hs_alloc_scratch(db, &s) != HS_SUCCESS)
    {
        state.SkipWithError("alloc scratch");
    }

    auto text = make_text(state.range(1));

    for (auto _ : state)
    {
        hs_stream_t *st = nullptr;

        if (hs_open_stream(db, 0, &st) != HS_SUCCESS)
        {
            state.SkipWithError("open stream failed");
        }

        auto data = text.data();

        for (auto i = 0; i < text.size(); i += page_size)
        {
            auto n = std::min(page_size, text.size() - i);

            if (hs_scan_stream(st, data + i, n, 0, s, on_match_event, nullptr) != HS_SUCCESS)
            {
                state.SkipWithError("scan failed");
            }
        }

        if (hs_close_stream(st, s, on_match_event, nullptr) != HS_SUCCESS)
        {
            state.SkipWithError("close stream failed");
        }
    }

    state.SetBytesProcessed(int64_t(state.iterations()) * int64_t(state.range(1)));

    hs_free_scratch(s);
}

BENCHMARK(BM_StreamScan)->ArgsProduct(args)->ArgNames({"regex", "size"});

#if HAVE_RE2

static void BM_RE2Match(benchmark::State &state)
{
    auto expr = benchData[benchCase(state.range(0))];

    RE2 pattern(expr);

    auto data = make_text(state.range(1));
    auto text = StringPiece(data.data(), data.size());

    for (auto _ : state)
    {
        if (RE2::FullMatch(text, pattern))
        {
            state.SkipWithError("scan failed");
        }
    }

    state.SetBytesProcessed(int64_t(state.iterations()) * int64_t(state.range(1)));
}

BENCHMARK(BM_RE2Match)->ArgsProduct(args)->ArgNames({"regex", "size"});

#endif

#if HAVE_PCRE2

static void BM_PCRE2Match(benchmark::State &state)
{
    auto expr = benchData[benchCase(state.range(0))];

    int err = 0;
    PCRE2_SIZE err_off = 0;

    pcre2_code *code = pcre2_compile(
        (PCRE2_SPTR)expr.c_str(),
        (PCRE2_SIZE)expr.size(),
        PCRE2_MULTILINE,
        &err,
        &err_off,
        nullptr);
    if (!code)
    {
        state.SkipWithError("compile failed");
    }

    pcre2_match_data *md = pcre2_match_data_create_from_pattern(code, nullptr);

    auto text = make_text(state.range(1));

    for (auto _ : state)
    {
        if (pcre2_match(code, (PCRE2_SPTR)text.data(), (PCRE2_SIZE)text.size(), 0, 0, md, nullptr) != PCRE2_ERROR_NOMATCH)
        {
            state.SkipWithError("scan failed");
        }
    }

    state.SetBytesProcessed(int64_t(state.iterations()) * int64_t(state.range(1)));

    pcre2_match_data_free(md);
    pcre2_code_free(code);
}

BENCHMARK(BM_PCRE2Match)->ArgsProduct(args)->ArgNames({"regex", "size"});

#endif

BENCHMARK_MAIN();
