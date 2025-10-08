import React from "react";
import "../styles/App.css";
import "../styles/Button.css";
import "../styles/Home.css";
import WoerterbuchPicture from "../images/woerterbuch.png";
import { useNavigate } from "react-router-dom";
import Title from "../components/Title";

const Home = () => {
    const navigate = useNavigate();

    return (
        <div className="page-container">
            <div className="page-frame">
                <Title/>
                <div className="main-menu-image">
                    <img src={WoerterbuchPicture} alt="logo" style={{width: "100%", height: "100%"}}/>
                </div>
                <div className="main-menu-buttons footer">
                    <div className="vocabulary-actions">
                        <span>
                            <button className="wb-button main-menu-button main-interface-button" onClick={() => navigate("/add")}>
                                Add a word
                            </button>
                        </span>
                        <span>
                            <button className="wb-button main-menu-button main-interface-button" onClick={() => navigate("/list")}>
                                Dictionary
                            </button>
                        </span>
                    </div>
                    <div className="train-button">
                        <button className="wb-button main-menu-button main-interface-button" onClick={() => navigate("/train")}>
                            Training mode
                        </button>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default Home;