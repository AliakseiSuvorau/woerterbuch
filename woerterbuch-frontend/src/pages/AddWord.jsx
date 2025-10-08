import React, {useState} from "react";
import BackButton from "../components/BackButton";
import {backendUrl} from '../constants/AppConstants';
import "../styles/App.css";
import Title from "../components/Title";
import "../styles/AddWord.css";
import tick from "../images/tick.svg";
import upload from "../images/upload.svg";

const AddWord = () => {
    const [selectedArticle, setSelectedArticle] = useState("der");
    const [word, setWord] = useState("");
    const [translation, setTranslation] = useState("")

    // Add a single word in web interface
    const handleAddWord = async () => {
        if (!selectedArticle || !word || !translation) {
            console.log(selectedArticle, word, translation)
            alert("Choose an article and enter a word and a translation!");
            return;
        }

        const response = await fetch(`${backendUrl}/dictionary/word/add`, {
            method: "POST",
            headers: {"Content-Type": "application/json"},
            body: JSON.stringify({article: selectedArticle, word, translation}),
        });

        if (response.ok) {
            setSelectedArticle("der");
            setWord("");
            setTranslation("");
            showToast("Word was added successfully!");
        } else {
            alert("Error while adding the word");
        }
    };

    // Upload dict from file
    let dict = null;  // Choose a file
    const handleUploadDictFromCsv = async () => {
        const input = document.createElement("input");
        input.type = "file";
        input.accept = ".csv";

        input.onchange = async (event) => {
            const file = event.target.files[0];
            if (!file) return;

            if (!file.name.toLowerCase().endsWith(".csv")) {
                alert("Please select a .csv file");
                return;
            }

            dict = await file.text();

            try {
                const response = await fetch(`${backendUrl}/dictionary/upload`, {
                    method: "POST",
                    headers: {"Content-Type": "text/plain"},
                    body: dict,
                });

                if (response.ok) {
                    showToast("Dictionary uploaded successfully!");
                } else {
                    alert("Error while uploading dictionary");
                }
            } catch (error) {
                console.error(error);
                alert("Network error while uploading dictionary");
            }
        };
        input.click();
    };

    return (
        <div className="page-container">
            <div className="page-frame">
                <Title/>
                <div className="page-subtitle">
                    <h2>Enter a word</h2>
                </div>
                <div className="new-word">
                    <select
                        value={selectedArticle}
                        onChange={(e) => setSelectedArticle(e.target.value)}>
                        <option value="der">der</option>
                        <option value="die">die</option>
                        <option value="das">das</option>
                    </select>
                    <div style={{display: "flex", flexDirection: "column"}}>
                        <input
                            type="text"
                            value={word}
                            onChange={(e) => setWord(e.target.value)}
                            placeholder="Enter a word"
                        />
                        <input
                            type="text"
                            value={translation}
                            onChange={(e) => setTranslation(e.target.value)}
                            placeholder="Enter a translation"
                        />
                    </div>
                    <button className="wb-button save-button"
                            onClick={handleAddWord}>
                        <img src={tick} alt={"Save"}/>
                    </button>
                </div>
                <div className="page-subtitle">
                    <h2>Or upload a csv file</h2>
                </div>
                <div className="upload-dict-from-csv">
                    <button className="wb-button upload-button"
                            onClick={handleUploadDictFromCsv}>
                        <img src={upload} alt={"Upload"}/>
                    </button>
                </div>
                <div className="footer">
                    <BackButton/>
                </div>
            </div>
        </div>
    );
};

export default AddWord;

function showToast(message) {
    const toast = document.createElement("div");
    toast.textContent = message;
    toast.style.position = "fixed";
    toast.style.bottom = "20px";
    toast.style.right = "20px";
    toast.style.backgroundColor = "rgba(0, 0, 0, 0.8)";
    toast.style.color = "#fff";
    toast.style.padding = "10px 16px";
    toast.style.borderRadius = "8px";
    toast.style.fontSize = "14px";
    toast.style.boxShadow = "0 2px 6px rgba(0,0,0,0.3)";
    toast.style.zIndex = "9999";
    toast.style.opacity = "0";
    toast.style.transition = "opacity 0.3s ease";

    document.body.appendChild(toast);

    requestAnimationFrame(() => {
        toast.style.opacity = "1";
    });

    setTimeout(() => {
        toast.style.opacity = "0";
        setTimeout(() => toast.remove(), 300);
    }, 2000);
}
