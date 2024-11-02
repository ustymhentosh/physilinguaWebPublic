const queryString = window.location.search;
const urlParams = new URLSearchParams(queryString);
const prenumber = urlParams.get('which');
const yellow = urlParams.get('yellow');

var mess = document.getElementById("message");
var butt = document.getElementById("main-button")

if (yellow === "true") {
    displayMessage("Цю задачу уже переклали, проте ще не перевірили. Розгляньте можливість перекласи іншу задачу.", "yellow")
}

function displayMessage(messageText, whichColor) {
    mess.innerHTML = messageText;
    mess.style.display = 'inline';
    if (whichColor === "red") {
        mess.style.backgroundColor = "#ff9494";
    } else if (whichColor === "orange") {
        mess.style.backgroundColor = "#ffd194"
    } else {
        mess.style.backgroundColor = "#fffb94"
    }
    document.getElementById("forma").style.paddingTop = "1rem";
}

function hideMessage() {
    mess.style.display = 'none';
    document.getElementById("forma").style.paddingTop = "2rem";
}

var renderedWindow = document.getElementById("renderedWindow");
renderedWindow.classList.add('no-text');
var rawUserText = document.getElementById("rawUserText");

rawUserText.oninput = function () {
    var renderedWindow = document.getElementById("renderedWindow");
    renderedWindow.innerHTML = rawUserText.value;
    MathJax.Hub.Queue(["Typeset", MathJax.Hub, "previewParagraph"]);
    if (rawUserText.value === "") {
        renderedWindow.innerHTML = "Тут відображатиметься ваша задача..."
        renderedWindow.classList.add('no-text');
    } else {
        renderedWindow.classList.remove('no-text');
    }
};

var renderedAnswerWindow = document.getElementById("renderedAnswerWindow");
renderedAnswerWindow.classList.add('no-text');
var rawUserAnswerText = document.getElementById("rawUserAnswerText");

rawUserAnswerText.oninput = function () {
    var renderedAnswerWindow = document.getElementById("renderedAnswerWindow");
    renderedAnswerWindow.innerHTML = rawUserAnswerText.value;
    MathJax.Hub.Queue(["Typeset", MathJax.Hub, "previewParagraph"]);
    if (rawUserAnswerText.value === "") {
        renderedAnswerWindow.innerHTML = "відповідь..."
        renderedAnswerWindow.classList.add('no-text');
    } else {
        renderedAnswerWindow.classList.remove('no-text');
    }
};

var fillerState = false

const COUNTRY_LIST = mydata;
const selectMenu = document.getElementById('selectMenu');
const selectBtn = document.getElementById('selectBtn');
const selectBtnSpan = document.querySelector('#selectBtn span');
const searchInput = document.getElementById('countryInput');
const countryList = document.getElementById('CountryList');
const hiddenInput = document.getElementById('hiddenInput');

let selectedCountry = false;

if (prenumber !== null) {
    selectBtnSpan.innerText = prenumber;
    selectBtn.style.borderColor = "#edb7f6";
    document.getElementById("number-error").style.display = "none"
    selectedCountry = prenumber;
    hiddenInput.value = prenumber;
    searchInput.value = '';
}

spanNumber = document.getElementById("selected-number")

function insertCountries(list) {
    let li = '';
    list.forEach((name) => {
        if (selectedCountry === name) {
            if (pending_problems.includes(name)) {
                li += `<li onclick="countryClick(this)" class="selected" style="color: red">${name}</li>`;
            } else if (ready_problems.includes(name)) {
                li += `<li onclick="countryClick(this)" class="selected" style="color: red;">${name}</li>`;
            } else {
                li += `<li onclick="countryClick(this)" class="selected">${name}</li>`;
            }
        }
        else {
            if (pending_problems.includes(name)) {
                li += `<li onclick="countryClick(this)" style="color: red">${name}</li>`;
            } else if (ready_problems.includes(name)) {
                li += `<li onclick="countryClick(this)" style="color: red;">${name}</li>`;
            } else {
                li += `<li onclick="countryClick(this)">${name}</li>`;
            }
        }
    });
    countryList.innerHTML = li;
}
insertCountries(COUNTRY_LIST);

selectBtn.onclick = function () {
    selectMenu.classList.toggle('active');
    if (selectMenu.classList.contains('active')) searchInput.focus();
    if (fillerState) {
        fillerState = false;
        document.getElementById("text-cont").style.marginLeft = "0rem";
    } else {
        fillerState = true;
        document.getElementById("text-cont").style.marginLeft = "9rem";
    }
}

function sendGetRequest(problemNum, callback) {
    var currentUrl = window.location.href;
    var checkUrl = currentUrl + "/check";
    var xhr = new XMLHttpRequest();

    xhr.open("GET", checkUrl + "?problemNum=" + problemNum, true);

    xhr.onload = function () {
        if (xhr.status >= 200 && xhr.status < 300) {
            var responseData = JSON.parse(xhr.responseText);
            callback(responseData.permission);
        } else {
            console.error('Request failed with status:', xhr.status);
            callback(false); // Callback with false in case of error
        }
    };
    xhr.onerror = function () {
        console.error('Request failed');
        callback(false); // Callback with false in case of error
    };
    xhr.send();
}

function countryClick(el) {
    selectBtnSpan.innerText = el.innerText;
    selectBtn.style.borderColor = "#edb7f6";
    document.getElementById("number-error").style.display = "none"
    selectedCountry = el.innerText;
    hiddenInput.value = el.innerText;
    selectMenu.classList.toggle('active');
    searchInput.value = '';
    insertCountries(COUNTRY_LIST);
    if (fillerState) {
        document.getElementById("text-cont").style.marginLeft = "0rem";
        fillerState = false;
    } else {
        document.getElementById("text-cont").style.marginLeft = "9rem";
        fillerState = true;
    }

    if (pending_problems.includes(selectedCountry)) {
        sendGetRequest(selectedCountry, function (permission) {
            if (permission) {
                displayMessage("Цю задачу уже переклали, проте ще не перевірили. Розгляньте можливість перекласи іншу задачу.", "yellow")
                document.getElementById("button-holder").appendChild(butt);
            } else {
                displayMessage("Цю задачу уже переклали кілька разів. Будь ласка оберіть іншу задачу.", "orange")
                butt.remove();
            }
        });
    } else if (ready_problems.includes(selectedCountry)) {
        displayMessage("Цю задачу уже переклали та опублікували. Розгляньте можливість перекласи іншу задачу.", "red")
        butt.remove()
    } else {
        hideMessage()
        document.getElementById("button-holder").appendChild(butt);
    }
}

const filteredCountries = function (e) {
    let keyword = searchInput.value.toLowerCase();
    let filteredResult = COUNTRY_LIST.filter((country) => {
        country = country.toLowerCase();
        return country.indexOf(keyword) > -1;
    });
    insertCountries(filteredResult);
}
searchInput.addEventListener('keyup', filteredCountries);


document.addEventListener('DOMContentLoaded', function () {
    var myForm = document.getElementById('forma');
    var popup = document.getElementById('popup');

    myForm.addEventListener('submit', function (event) {
        event.preventDefault();
        handleFormSubmission();
    });

    function handleFormSubmission() {
        // var requiredFields = myForm.querySelectorAll('[required]');
        // requiredFields.forEach((field) => {
        //     if (field.value.trim() === '') {
        //       const fieldName = field.getAttribute('name');
        //       const errorMessage = document.createElement('p');
        //       errorMessage.textContent = `${fieldName} is required.`;
        //       errorMessagesDiv.appendChild(errorMessage);
        //     }});
        var number = document.getElementById('hiddenInput').value;
        var formData = new FormData(myForm);

        if (number) {
            fetch('./formsubmit', {
                method: 'POST',
                body: formData
            })
                .then(response => response.json())
                .then(data => {
                    console.log('Success:', data);
                    myForm.reset();
                    showPopup();
                })
                .catch(error => {
                    console.error('Error:', error);
                });
        } else {
            document.getElementById("number-error").style.display = "block";
            console.log('Form is not valid. Please fill in all required fields.');
        }
    }

    function showPopup() {
        renderedAnswerWindow.innerHTML = "надіслано..."
        renderedWindow.innerHTML = "надіслано..."
        popup.style.display = 'block'; // Display the pop-up
        setTimeout(function () {
            location.reload()
            popup.style.display = 'none';
        }, 3000);
    }
});