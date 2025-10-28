import { Link, useNavigate } from "react-router-dom";

export default function NavBar({ authenticated, username }: { authenticated: boolean, username: string }) {
  const navigate = useNavigate();

  const logout = () => {
    localStorage.removeItem("token");
    localStorage.removeItem("username");
    navigate("/login");
  };

  return (
    <nav className="bg-blue-100 text-blue-900 flex justify-between items-center px-6 py-3 shadow">
      <h1 className="font-bold text-lg">Go-u-sh</h1>
      {authenticated ? (
        <div className="flex items-right gap-4">
          <h2 className="font-semibold">Welcome {username}</h2>
          <button
            onClick={logout}
            className="bg-red-500 text-white px-3 py-1 rounded hover:bg-red-600"
          >
            Logout
          </button>
        </div>
      ) : (
        <div className="space-x-4">
          <Link to="/login" className="hover:underline">
            Login
          </Link>
          <Link to="/register" className="hover:underline">
            Register
          </Link>
        </div>
      )}
    </nav>
  )
}