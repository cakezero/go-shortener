import { useState } from "react";
import { useNavigate, Link } from "react-router-dom";
import { makeRequest } from "../axios";

export default function Login() {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const navigate = useNavigate();

  const handleLogin = async () => {
    const { data } = await makeRequest({
      endpoint: "/auth/login",
      dataOrQuery: { username, password },
      method: "POST"
    });

    console.log({ data })

    localStorage.setItem("token", data.token);
    localStorage.setItem("username", data.user.username);
    localStorage.setItem("id", data.user._id);

    navigate("/");
  };

  return (
    <div className="min-h-screen bg-white flex items-center justify-center">
      <form
        className="p-6 bg-blue-100 rounded-lg shadow-md w-80"
      >
        <h2 className="text-xl font-semibold mb-4">Login</h2>
        <input
          type="text"
          placeholder="Username"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
          className="border px-3 py-2 mb-4 w-full rounded"
        />
        <input
          type="password"
          placeholder="Password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          className="border px-3 py-2 mb-4 w-full rounded"
        />
        <button type="button" onClick={handleLogin} className="bg-blue-500 text-white px-4 py-2 rounded w-full">
          Login
        </button>
        <p className="mt-3 text-sm text-center">
          No account?{" "}
          <Link to="/register" className="text-blue-700 hover:underline">
            Register
          </Link>
        </p>
      </form>
    </div>
  );
}
