import React, { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import BackButton from "../components/BackButton";
import { backendUrl, defaultPageSize } from "../constants/AppConstants";
import "../styles/App.css";
import "../styles/ListWords.css";
import del from "../images/delete.svg";
import edit from "../images/edit.svg";
import add from "../images/add.svg";
import left from "../images/leftArrow.svg";
import right from "../images/rightArrow.svg";
import tick from "../images/tick.svg";
import Title from "../components/Title";

const ListWords = () => {
    const [words, setWords] = useState([]);
    const [page, setPage] = useState(1);
    const [pageSize] = useState(defaultPageSize);
    const navigate = useNavigate();

    useEffect(() => {
        fetch(`${backendUrl}/dictionary/list?page=${page}&size=${pageSize}`)
            .then((res) => res.json())
            .then((data) => {
                if (!data) {
                    setWords([])
                } else {
                    setWords(data)
                }
            })
            .catch(() => alert("Error while fetching words"));
    }, [page]);

    const handleDelete = async (id) => {
        try {
            const response = await fetch(`${backendUrl}/dictionary/word/delete`,  {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ id })
            });

            if (!response.ok) {
                alert("Error deleting a word");
                return;
            }

            await handleReloadPage()
        } catch (error) {
            console.error(error);
            alert("Error deleting a word")
        }
    };

    const handleEdit = (id, currentArticle, currentWord, currentTranslation) => {
        setWords(words.map(word =>
            word.id === id
                ? { ...word, isEditing: true, newArticle: currentArticle, newWord: currentWord, newTranslation: currentTranslation }
                : word
        ));
    };

    const handleSave = (id, newArticle, newWord, newTranslation) => {
        fetch(`${backendUrl}/dictionary/word/edit`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ id, article: newArticle, word: newWord, translation: newTranslation })
        })
            .then(() => {
                setWords(words.map(word =>
                    word.id === id
                        ? { ...word, isEditing: false, article: newArticle, word: newWord, translation: newTranslation }
                        : word
                ));
            })
            .catch(() => alert("Error editing a word"));
    };

    const handleNextPage = () => {
        fetch(`${backendUrl}/dictionary/list?page=${page + 1}&size=${pageSize}`)
            .then((res) => res.json())
            .then((data) => {
                if (data && data.length > 0) {
                    setPage(prev => prev + 1);
                }
            })
            .catch(() => alert("Error checking next page"));
    };

    const handleReloadPage = async () => {
        try {
            const res = await fetch(`${backendUrl}/dictionary/list?page=${page}&size=${pageSize}`);
            const data = await res.json();
            if (data) {
                setWords(data);
            } else {
                if (page > 1) {
                    const res = await fetch(`${backendUrl}/dictionary/list?page=${page-1}&size=${pageSize}`);
                    const data = await res.json();
                    setWords(data)
                    setPage(prev => prev-1)
                } else {
                    setWords([])
                }
            }
        } catch(error) {
            alert("Error reloading page")
        }
    };

    return (
        <div className="page-container">
            <div className="page-frame">
                <Title/>
                <div className="page-subtitle">
                    <h2>Dictionary</h2>
                </div>
                <div className="list-container">
                    <div>
                        {words.map((item) => (
                            <>
                                {item.isEditing ? (
                                    <div className="edit-word">
                                        <select
                                            value={item.newArticle}
                                            onChange={(e) =>
                                                setWords(words.map(word =>
                                                    word.id === item.id
                                                        ? {...word, newArticle: e.target.value}
                                                        : word
                                                ))
                                            }
                                            style={{width: "fit-content", aspectRatio: "1 / 1"}}
                                        >
                                            <option value="der">der</option>
                                            <option value="die">die</option>
                                            <option value="das">das</option>
                                        </select>
                                        <div style={{display: "flex", flexDirection: "column"}}>
                                            <input
                                                type="text"
                                                value={item.newWord}
                                                onChange={(e) =>
                                                    setWords(words.map(word =>
                                                        word.id === item.id
                                                            ? {...word, newWord: e.target.value}
                                                            : word
                                                    ))
                                                }
                                            />
                                            <input
                                                type="text"
                                                value={item.newTranslation}
                                                onChange={(e) =>
                                                    setWords(words.map(word =>
                                                        word.id === item.id
                                                            ? {...word, newTranslation: e.target.value}
                                                            : word
                                                    ))
                                                }
                                            />
                                        </div>
                                        <button className="wb-button save-button"
                                                onClick={() => handleSave(item.id, item.newArticle, item.newWord, item.newTranslation)}>
                                            <img src={tick} alt={"Save"}/>
                                        </button>
                                    </div>
                                ) : (
                                    <div className="list-item">
                                        <div className="word">
                                            <span>
                                                <p style={{margin: 0}}>{item.article} {item.word}</ p>
                                            </span>
                                            <span>
                                                <p style={{
                                                    margin: 0,
                                                    color: "#7c7a7a",
                                                    fontSize: "small"
                                                }}>{item.translation}</p>
                                            </span>
                                        </div>
                                        <button className="wb-button edit-delete-button"
                                                onClick={() => handleEdit(item.id, item.article, item.word, item.translation)}>
                                            <img src={edit} alt="Edit"/>
                                        </button>
                                        <button className="wb-button edit-delete-button"
                                                onClick={() => handleDelete(item.id)}>
                                            <img src={del} alt="Delete"/>
                                        </button>
                                    </div>
                                )}
                            </>
                        ))}
                    </div>
                    <div style={{display: "flex", justifyContent: "center"}}>
                        <button className="wb-button add-button" onClick={() => navigate("/add")}>
                            <img src={add} alt="Add"/>
                        </button>
                    </div>
                </div>
                <div className="dict-page-footer footer">
                    <div className="pagination">
                        <span>
                            <button className="wb-button change-page-button main-interface-button"
                                    onClick={() => setPage(prev => Math.max(prev - 1, 1))}
                                    style={{display: "flex", justifyContent: "center"}}>
                                <img src={left} alt="Previous"/>
                            </button>
                        </span>
                        <span className="page-number">Page {page}</span>
                        <span>
                            <button className="wb-button change-page-button main-interface-button"
                                    onClick={handleNextPage}
                                    style={{display: "flex", justifyContent: "center"}}>
                                 <img src={right} alt="Next"/>
                            </button>
                        </span>
                    </div>
                    <BackButton/>
                </div>
            </div>
        </div>
    );
};

export default ListWords;
