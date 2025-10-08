import React, { useState } from "react";
import BackButton from "../components/BackButton";
import {backendUrl} from '../constants/AppConstants';
import "../styles/App.css";
import Title from "../components/Title";
import "../styles/AddWord.css";
import tick from "../images/tick.svg";

const AddWord = () => {
  const [selectedArticle, setSelectedArticle] = useState("der");
  const [word, setWord] = useState("");
  const [translation, setTranslation] = useState("")

  const handleAddWord = async () => {
    if (!selectedArticle || !word || !translation) {
        console.log(selectedArticle, word, translation)
      alert("Choose an article and enter a word and a translation!");
      return;
    }

    const response = await fetch(`${backendUrl}/dictionary/word/add`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ article: selectedArticle, word, translation }),
    });

    if (response.ok) {
      setSelectedArticle("der");
      setWord("");
      setTranslation("");
    } else {
      alert("Error while adding the word");
    }
  };

  return (
      <div className="page-container">
          <div className="page-frame">
              <Title/>
              <div className="page-subtitle">
                  <h2>Add a word</h2>
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
              <div className="footer">
                  <BackButton/>
              </div>
          </div>
      </div>
  );
};

export default AddWord;
