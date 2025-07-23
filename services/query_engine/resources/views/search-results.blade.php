<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    @vite('resources/css/app.css')
    <title>Search Results for "{{ $originalQuery }}"</title>
</head>

<body>
    <x-search-bar />
    <div class="results-counter">
        @php
            $suggestion = '';
            if ($suggestions) {
                $suggestion = "Did you mean $query? ";
            }
        @endphp
        <span id="suggestion">{{ $suggestion }}</span><span>Showing {{ $total }} results for {{ $query }}</span>
        @if ($suggestions)
            <p>Search for <a id="suggestion-link"
                    href="{{ route('search', ['processed_query' => $originalQuery]) }}">{{ $originalQuery }}</a>
                instead
            </p>
        @endif
    </div>

    <div class="results-container">
        <div class="pages-container col-span-2">
            @if (count($results) > 0)
                <ul>
                    @foreach ($results as $res)
                        <x-search-result url="{{ $res->_id }}" title="{{ $res->title }}"
                        text="{{ $res->summary_text }}" />
                    @endforeach
                </ul>
            @else
                <p>No results found</p>
            @endif
        </div>
        @if (@page != null && $page == 1)
            <div class="images-container flex flex-col items-center">
                <h2 class="top-images-text"> {{ ucwords($query) }} </h2>
                @foreach ($topImages as $res)
                    <div class="my-2">
                        <x-image-container url="{{ $res->_id }}" alt="{{ $res->alt }}"
                            title="{{ $res->page_title }}" page="{{ $res->page_url }}"
                            text="{{ $res->page_text }}" />
                    </div>
                @endforeach
            </div>
        @endif
        <div class="flex flex-col justify-center items-center">
            <x-pagination-bar totalResults="{{ $total }}" />
        </div>
    </div>
     <footer>
        <p> <a href="https://github.com/kawilkinson/superseek">Support/Contribute to the project!</a></p>
        <p>
            <a href="https://github.com/kawilkinson">Github</a>
            <a href="https://www.linkedin.com/in/kyleawilkinson/">LinkedIn</a>
        </p>
        <p id="copyright">
            &copy; {{ now()->year }} All rights reserved.
        </p>
     </footer>
</body>

</html>