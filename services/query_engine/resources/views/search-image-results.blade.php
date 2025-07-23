<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    @vite('resources/css/app.css')
    <title>Search Image Results for "{{ $query }}"</title>
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
        <span id="suggestion">{{ $suggestion }}</span><span>Showing {{ $total }} results for 
            {{ $query }}</span>
        @if ($suggestions)
            <p>Search for <a id="suggestion-link"
                href="{{ route('search-images', ['processed_query' => $originalQuery]) }}">{{ $originalQuery }}</a>
                instead
            </p>
        @endif
    </div>
    <div class="results-images-container">
        @if (count($results) > 0)
            <ul class="flex flex-wrap gap-4">
                @foreach ($results as $res)
                    <li class="m-2 flex-shrink-0 w-1/6">
                        <x-image-container url="{{ $res->_id }}" alt="{{ $res->alt }}"
                            title="{{ $res->page_title }}" page="{{ $res->page_url }}" text="{{ $res->page_text }}" />
                    </li>
                @endforeach
            </ul>
        @else
            <p>No results found</p>
        @endif
    </div>
    <div class="flex flex-col justify-center items-center">
        <x-pagination-bar totalResults="{{ $total }}" />
    </div>
    <!-- Footer -->
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