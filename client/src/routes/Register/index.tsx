import { useState } from "react";
import "./index.css";
import { login, signUp } from "../../api/signup-login/api";
import { useNavigate } from "react-router";

const Register = () => {
  const [usernameInput, setUsernameInput] = useState("");
  const [passwordInput, setPasswordInput] = useState("");
  const [confirmPasswordInput, setConfirmPasswordInput] = useState("");
  const [isSigningUp, setIsSigningUp] = useState(true);
  const navigate = useNavigate();

  const handleRegister = async () => {
    try {
      if (isSigningUp) {
        await signUp(usernameInput, passwordInput);
      }
      const data = await login(usernameInput, passwordInput);
      console.log("data : ", data);

      localStorage.setItem("token", data.token);
      localStorage.setItem("name", data.name);
      setUsernameInput("");
      setPasswordInput("");
      navigate("/todos");
    } catch (err) {
      console.log("merde");
      console.log(err);
    }
  };

  return (
    <main>
      <h1>{isSigningUp ? "Signup" : "Login"}</h1>
      <form
        onSubmit={async (e) => {
          e.preventDefault();
          if (passwordInput !== confirmPasswordInput) {
            console.log("passwords don't match");
            return;
          }
          await handleRegister();
        }}
      >
        <div className="form-component">
          <label htmlFor="username">username : </label>
          <input
            type="text"
            id="username"
            placeholder="your username.."
            value={usernameInput}
            onChange={(e) => {
              setUsernameInput(e.currentTarget.value);
            }}
          />
        </div>
        <div className="form-component">
          <label htmlFor="password">password : </label>
          <input
            type="password"
            id="password"
            placeholder="your password.."
            value={passwordInput}
            onChange={(e) => {
              setPasswordInput(e.currentTarget.value);
            }}
          />
        </div>
        {isSigningUp ? (
          <div className="form-component">
            <label htmlFor="confirm-password">confirm password : </label>
            <input
              type="password"
              id="confirm-password"
              placeholder="confirm password.."
              value={confirmPasswordInput}
              onChange={(e) => {
                setConfirmPasswordInput(e.currentTarget.value);
              }}
            />
          </div>
        ) : undefined}
        <button type="submit">Submit</button>
      </form>
      <a
        style={{
          marginTop: "20px",
          cursor: "pointer",
          textDecoration: "underline",
        }}
        onClick={() => {
          setIsSigningUp(!isSigningUp);
        }}
      >
        {isSigningUp ? "Already have an account ?" : "Create an account ?"}
      </a>
    </main>
  );
};

export default Register;
