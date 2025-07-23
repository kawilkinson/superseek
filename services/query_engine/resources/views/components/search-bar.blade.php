<div class="search-bar-container">
    <div class="top-part">
        @php
            $currentPage = request()->query('page', 1);
            $currentQuery = request()->query('q');
            $currentPath = request()->query();
            $currentSearchFunction = explode('/', $currentPath)[1];
        @endphp
        <a href="https://superseek.app/" class="logo-container">
            Superseek
        </a>
        @php
            $currentAction = '/api/search/';
            if ($currentSearchFunction == 'searchImages') {
                $currentAction = '/api/search-images';
            }
        @endphp
        <form action="{{ $currentAction }}" method="GET" class="px-3">
            <input type="text" id="search-bar" name="q" placeholder="Search..." autocomplete="off" />
            <div id="search-suggestions" class="absolute bg-white search-results hidden">
                <ul id="search-suggestions-list"></ul>
            </div>
            <button type="submit" class="btn" id="search-button">
                Search it!
            </button>
            <button type="button" class="btn" onclick="window.location.href='/api/random'">
                Random Page!
            </button>
        </form>
    </div>
    <div class="bottom-part">
        @php
            $searchQuery = $currentQuery;
            if (!$searchQuery) {
                $searchQuery = request()->query('processed_query');
            }

            $imagesUrl = '/api/search-images?q=' . $searchQuery;
            $pagesUrl = '/api/search?q=' . $searchQuery;

            $imagesActive = '';
            $pagesActive = '';
            if ($currentSearchFunction == 'search') {
                $imagesActive = '';
                $pagesActive = 'active';
            } else {
                $imagesActive = 'active';
                pagesActive = '';
            }
        @endphp
        <a href="{{ $pagesUrl }}" class="tab {{ $pagesActive }}">
            <span> PAGES </span>
        </a>
        <a href="{{ $imagesUrl }}" class="tab {{ $imagesActive }}">
            <span> IMAGES </span>
        </a>
    </div>
</div>

<script>
    document.addEventListener('DOMContentLoaded', function() {
        const searchBar = document.getElementById('search-bar');
        const resultsContainer = document.getElementById('search-suggestions');
        const resultsList = document.getElementById('search-suggestions-list');
        let debounceTimer;

        function debounce(func, delay) {
            clearTimeout(debounceTimer);
            debounceTimer = setTimeout(func, delay);
        }

        searchBar.addEventListener('input', function() {
            debounce(() => {
                if (searchBar.value.length >= 2) {
                    fetchSuggestions(searchBar.value);
                } else {
                    fetchTopSearches();
                }
            }, 300);
        });

        searchBar.addEventListener('focus', function() {
            if (searchBar.value.length < 2) {
                fetchTopSearches();
            } else {
                fetchSuggestions(searchBar.value);
            }
        });

        // for hiding results when clicking outside of search bar
        document.addEventListener('click', function(event) {
            if (!resultsContainer.contains(event.target) && event.target !== searchBar) {
                hideResults();
            }
        });

        // for hiding results when pressing escape
        document.addEventListener('keydown', function(event) {
            if (event.key === 'Escape') {
                hideResults();
                searchBar.blur();
            }
        });

        function hideResults() {
            resultsContainer.classList.add('hidden');
        }

        function fetchSuggestions(query) {
            fetch(`{{ route('get-search-suggestions') }}?q=${encodeURIComponent(query)}`)
                .then(response => {
                    if (!response.ok) {
                        throw new Error('Failed to fetch search results');
                    }
                    return response.json();
                })
                .then(data => {
                    displayResults(data.searches);
                })
                .catch(error => {
                    hideResults();
                    console.error(error);
                });
        }

        function displayResults(results) {
            resultsList.innerHTML = '';

            results.forEach(result => {
                const li = document.createElement('li');
                li.className = 'search-suggestion';
                li.textContent = result;

                li.addEventListener('click', function() {
                    searchBar.value = result;
                    hideResults();
                });

                resultsList.appendChild(li);
            });

            resultsContainer.classList.remove('hidden');
        }
    });
</script>