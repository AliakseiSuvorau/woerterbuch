import React, { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import "../styles/Button.css";
import "../styles/App.css";
import "../styles/Train.css";
import {articles, backendUrl} from "../constants/AppConstants";
import BackButton from "../components/BackButton";
import Title from "../components/Title";

const Train = () => {
    const [words, setWords] = useState([]);
    const [currentIndex, setCurrentIndex] = useState(0);
    const [selectedArticle, setSelectedArticle] = useState(null);
    const [correctCount, setCorrectCount] = useState(0);
    const [disabledArticles, setDisabledArticles] = useState([]);
    const [feedback, setFeedback] = useState("");

    useEffect(() => {
        fetchWords();
    }, []);

    const fetchWords = async () => {
        try {
            const response = await fetch(`${backendUrl}/dictionary/getRandomBatch`);
            const data = await response.json();
            setWords(data);
            setCurrentIndex(0);
            setSelectedArticle(null);
            setDisabledArticles([]);
            setFeedback("");
        } catch (error) {
            console.error("Error while fetching words: ", error);
        }
    };

    const handleArticleClick = (article) => {
        if (!words.length) return;

        const correctArticle = words[currentIndex].article;
        if (article === correctArticle) {
            setFeedback("Correct!");
            setCorrectCount(correctCount + 1);
            setTimeout(() => {
                if (currentIndex + 1 < words.length) {
                    setCurrentIndex(currentIndex + 1);
                    setSelectedArticle(null);
                    setDisabledArticles([]);
                    setFeedback("");
                } else {
                    fetchWords();
                }
            }, 400);
        } else {
            setFeedback("Try again!");
            setDisabledArticles([...disabledArticles, article]);
            setCorrectCount(0);
        }
    };

    return (
        <div className="page-container">
            <div className="page-frame">
                <Title/>
                <div className="train-container">
                    <div className="page-subtitle">
                        <h2>Training mode</h2>
                    </div>
                    <div className="train-counter">
                        <p>Correct answers in a row: {correctCount}</p>
                    </div>
                    {words.length > 0 && (
                        <div className="word-with-translation-and-choice-buttons">
                            <div className="word-with-translation">
                                <h2 style={{margin: 0}}>{words[currentIndex].word}</h2>
                                <p style={{margin: 0, color: "#7c7a7a"}}>{words[currentIndex].translation}</p>
                            </div>
                            <div className="article-buttons">
                                {articles.map((article) => (
                                    <button
                                        key={article}
                                        className={`wb-button article-button ${selectedArticle === article ? "active-button" : ""}`}
                                        onClick={() => handleArticleClick(article)}
                                        disabled={disabledArticles.includes(article)}
                                    >
                                        {article}
                                    </button>
                                ))}
                            </div>
                        </div>
                    )}
                    <div className="feedback">
                        <p>{feedback}</p>
                    </div>
                </div>
                <div className="footer">
                    <BackButton/>
                </div>
            </div>
        </div>
    );
};

export default Train;