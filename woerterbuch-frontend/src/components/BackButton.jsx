import React from "react";
import { useNavigate } from "react-router-dom";
import "../styles/Common.css";

const BackButton = () => {
  const navigate = useNavigate();
  return (
      <div className="back-button-footer">
          <button className="wb-button main-interface-button" onClick={() => navigate("/")}>
              Back
          </button>
      </div>
  );
};

export default BackButton;
