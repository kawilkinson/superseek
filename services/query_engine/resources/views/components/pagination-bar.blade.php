<div class="pagination-bar-container">
    <div class="flex gap-1 m-5">
        @php
            $osToRender = max(min($totalPages, 10), 2); // clamp down max pages
            $currentPage = request()->query('page', 1);
            $currentQuery = request()->query('q');
            $currentPath = request()->path();
        @endphp
        <span class="search-letter">S</span>
        @for ($i = 1; $i <= $osToRender; $i++)
            @php
                $activeClass = $i == $currentPage ? 'active' : 'inactive';
                $url = url($currentPath) . '?q=' . $currentQuery . '&page=' . $i;
            @endphp
            <a href="{{ $url }}" class="search-page {{ $activeClass }}">
                <span class="search-letter">u</span>
                <p>{{ $i }} </p>
            </a>
        @endfor
        <span class="search-letter">p</span>
        <span class="search-letter">e</span>
        <span class="search-letter">r</span>
        <span class="search-letter">s</span>
        <span class="search-letter">e</span>
        <span class="search-letter">e</span>
        <span class="search-letter">k</span>
    </div>
</div>