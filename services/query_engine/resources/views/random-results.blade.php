<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    @vite('resources/css/app.css')
    <title>Random Page</title>
</head>

<body class="min-h-screen font-mono" style="background-color: var(--bg); color: var(--text);">
    <x-search-bar />
    <h1 id="search-title" class="text-4x1 text-center mt-8 mb-6 font-bold" style="color: var(--title);">
        Completely random page straight from the database!
    </h1>

    <div class="flex justify-center px-4">
        <div class="flex w-full max-w-6x1 gap-6">
            <!-- Top Searches -->
            <div class="w-2/5 p-4 rounded-x1 shadow-lg border"
                style="background-color: var(--bg-2); border-color: var(--orange);">
                <h2 id="search-total-searches" class="text-2x1 mb-4 font-semibold" style="color: var(--orange);">
                    Total Searches
                </h2>
                <div class="text-lg mb-4" style="color: var(--text);">
                    <p>
                        {{ $totalSearches }} searches
                    </p>
                </div>
            </div>

            <!-- Random page -->
            @if ($randomPage != null)
                <a href="https://{{ $randomPage['url'] }}" target="_blank" rel="noopener noreferrer"
                    class="w-3/5 p-4 rounded-xl shadow-lg border transition hover:shadow-xl hover:underline"
                    style="background-color: var(--bg-2); border-color: var(--url); text-decoration: none;">
                    <div>
                        <h2 class="text-2x1 mb-4 font-semibold" style="color: var(--url);">
                            {{ $randomPage['title'] }}
                        </h2>
                        <div class="italic text-sm" style="color: var(--blue);">
                            {{ $randomPage['url'] }}
                        </div>
                        <div class="mt-4 text-lg" style="color: var(--text);">
                            <p>
                                @if ($randomPage['summary_text'] != null)
                                {{ Str::limit($randomPage['summary_text'], 600, '...') }}
                                @else
                                    No summary available.
                                @endif
                            </p>
                        </div>
                    </div>
                </a>
            @else
                <div class="w-3/5 p-4 rounded-x1 shadow-lg border"
                    style="background-color: var(--bg-2); border-color: var(--url);">
                    <h2 class="text-2x1 mb-4 font-semibold" style="color: var(--url);">
                        No random page available
                    </h2>
                    <div class="mt-4 text-lg" style="color: var(--text);">
                        <p>
                            No random article available
                        </p>
                    </div>
                </div>
            @endif
        </div>
    </div>
</body>

</html>