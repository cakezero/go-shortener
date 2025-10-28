import { useState } from "react";
import { useNavigate, Link } from "react-router-dom";
import { makeRequest } from "../axios";

export default function Register() {
  const [username, setUsername] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const navigate = useNavigate();

  const handleRegister = async () => {
    const { data } = await makeRequest({
      endpoint: "/auth/register",
      dataOrQuery: { username, email, password },
      method: "POST"
    });

    console.log({ data });

    localStorage.setItem("token", data.token);
    localStorage.setItem("id", data.user._id);
    localStorage.setItem("username", data.user.username);

    navigate("/");
  };

  return (
    <div className="min-h-screen bg-white flex items-center justify-center">
      <form
        className="p-6 bg-blue-100 rounded-lg shadow-md w-80"
      >
        <h2 className="text-xl font-semibold mb-4">Register</h2>
        <input
          type="text"
          placeholder="Username"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
          className="border px-3 py-2 mb-4 w-full rounded"
        />
        <input
          type="email"
          placeholder="Email"
          value={username}
          onChange={(e) => setEmail(e.target.value)}
          className="border px-3 py-2 mb-4 w-full rounded"
        />
        <input
          type="password"
          placeholder="Password"
          value={username}
          onChange={(e) => setPassword(e.target.value)}
          className="border px-3 py-2 mb-4 w-full rounded"
        />
        <button type="button" onClick={handleRegister} className="bg-blue-500 text-white px-4 py-2 rounded w-full">
          Register
        </button>
        <p className="mt-3 text-sm text-center">
          Already have an account?{" "}
          <Link to="/login" className="text-blue-700 hover:underline">
            Login
          </Link>
        </p>
      </form>
    </div>
  );
}
