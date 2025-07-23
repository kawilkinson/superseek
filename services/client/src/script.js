// const backendURL = `https://api.superseek.app/api`;
const backendURL = `http://127.0.0.1:8000/api`;

document.addEventListener("DOMContentLoaded", () => {
    const searchButton = document.getElementById("search-button");
    const randomButton = document.getElementById("random-button");
    const searchBar = document.getElementById("search-bar");
    let infoContainer;
    let formattedInfo;
    let html;

    if (searchButton) {
        searchButton.addEventListener("click", () => {
            const query = searchBar.value.trim();
            if (query) {
                search(query);
            } else {
                alert("Please enter a search query");
            }
        });
    }

    if (randomButton) {
        randomButton.addEventListener("click", () => {
            random();
        });
    }

    if (searchBar) {
        searchBar.value = "";
        searchBar.placeholder = "Type here to begin a search!";
        searchBar.addEventListener("keydown", (event) => {
            if (event.key == "Enter") {
                const query = searchBar.value.trim();
                if (query) {
                    search(query);
                } else {
                    alert("Please enter a search query");
                }
            }
        });
    }

    console.log(`Pinging backend at ${backendURL}...`);

    fetch(`$(backendURL)/stats`)
        .then((res) => res.json())
        .then((data) => {
            if (data.status === "up") {
                infoContainer = document.getElementById("info-container")

                formattedInfo = getFormattedInfoString(data.pages);

                html = `<span class="info">${formattedInfo}</span>`;
                infoContainer.innerHTML = html;
            } else if (data.status === "down") {
                serverDown();
            }
        })
        .catch((error) => {
            console.error("Error getting data:", error);
            serverDown();
        });

    function serverDown() {
        infoContainer = document.getElementById("info-container");
        infoContainer.innerHTML = `Server down - please come back later`;
        searchBar.disabled = true;
        searchBar.placeholder = "hello";
        searchBar.value = "Please come back later...";
        searchButton.disabled = true;
        randomButton.disabled = true;
    }
});

function getFormattedInfoString(pageCount) {
    const formatter = new Intl.NumberFormat("en", { notation: "compact" });

    let formattedString = "Contains ";
    const approxCount = formatter.format(pageCount);

    let order = approxCount.slice(-1);
    switch (order) {
        case "K":
            formattedString += `~ ${approxCount.slice(0, -1)} thousand`;
            break
        case "M":
            formattedString += `~ ${approxCount.slice(0, -1)} million`;
            break
        case "B":
            formattedString += `~ ${approxCount.slice(0, -1)} billion`;
            break
        default:
            formattedString += `${approxCount}`;
            break
    }

    formattedString += " results"

    return formattedString
}

async function search(query) {
    try {
        const encodedQuery = encodeURIComponent(query);
        const requestUrl = `${backendURL}/search?q=${encodedQuery}`;
        
        console.log(requestUrl);

        window.location.href = requestUrl;
    } catch (error) {
        console.log(error.message);
    }
}

async function random() {
    try {
        // const randomUrl = `${backendURL}/random`;
        const randomUrl = `http://localhost:8000/api/random`;

        window.location.href = randomUrl;
    } catch (error) {
        console.log(error.message);
    }
}