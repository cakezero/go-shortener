import { Trash2, Edit, X } from 'lucide-react';
import { useState, useEffect } from 'react';
import NavBar from '../components/NavBar';
import { makeRequest } from '../axios';

interface Url {
  id: string;
  longurl: string;
  shorturl: string;
}

export default function Body() {
  const [selectedIds, setSelectedIds] = useState<string[]>([]);
  const [inputValue, setInputValue] = useState('');
  const [authenticated, setAuthenticated] = useState(false);
  const [username, setUsername] = useState('');
  const [userId, setUserId] = useState<string | null>('');
  const [jwtToken, setJwtToken] = useState<string | null>('');
  const [urls, setUrls] = useState<Url[]>([]);
  const [editUrl, setEditUrl] = useState<Url | null>(null);

  useEffect(() => {
    const token = localStorage.getItem('token');
    const user = localStorage.getItem('username');
    const id = localStorage.getItem('id');

    if (token && user) {
      setAuthenticated(true);
      setUsername(user);
      setUserId(id);
      setJwtToken(token)
    }

    if (token) {
      (async () => {
        const { urls: urlsFetched } = await makeRequest({
          endpoint: "/fetch-urls",
          dataOrQuery: id!,
          method: "GET",
          jwtToken
        });
        if (urlsFetched.length != 0) {
          setUrls(urlsFetched);
        }
      })()
    }
  });

  // Select Url
  const handleSelectRow = (id: string) => {
    setSelectedIds((prev) =>
      prev.includes(id) ? prev.filter((sid) => sid !== id) : [...prev, id]
    );
  };

  // Delete Url
  const handleDelete = async (id: string) => {
    if (window.confirm('Are you sure you want to delete this url?')) {
      await makeRequest({
        endpoint: "/delete-url",
        dataOrQuery: id,
        method: "DELETE",
        jwtToken
      });
    }
  };

  const handleDeleteSelectedUrls = async () => {
    if (window.confirm('Delete selected urls?')) {
      await makeRequest({
        endpoint: "/delete-selected-urls",
        dataOrQuery: { ids: selectedIds },
        method: "DELETE",
        jwtToken
      })
    }
  };

  const handleDeleteAllUrls = async () => {
    if (window.confirm('Delete all urls?')) {
      await makeRequest({
        endpoint: "/delete-urls",
        dataOrQuery: userId,
        method: "DELETE",
        jwtToken
      })
    }
  };

  // Shorten Url
  const handleSubmit = async () => {
    if (!inputValue.trim()) return;
    const { data } = await makeRequest({
      endpoint: "/shorten",
      dataOrQuery: { url: inputValue, id: userId },
      method: "POST"
    });

    if (data) {
      setUrls([...urls, data]);
    }

    setInputValue('');
  };

  // Edit long Url
  const handleSaveEdit = () => {
    if (!editUrl) return;
    setUrls(urls.map((i) => (i.id === editUrl.id ? editUrl : i)));
    setEditUrl(null);
  };

  const handleEdit = (url: Url) => {
    setEditUrl(url);
  };

  return (
    <div className="min-h-screen bg-white text-black flex flex-col">
      <NavBar authenticated={authenticated} username={username} />
      <main className="flex flex-col items-center justify-start mt-10 w-full">
        <form className="flex space-x-3 mb-6">
          <input
            value={inputValue}
            onChange={(e) => setInputValue(e.target.value)}
            className="border border-blue-300 rounded px-3 py-2"
            placeholder="Enter URL with its protocol/host e.g https://example.com or http://example.com"
          />
          <button onClick={handleSubmit} type="button" className="bg-blue-500 text-white px-4 py-2 rounded">Shorten</button>
        </form>

        <div className="relative w-3/4">
          {authenticated && urls.length > 0 && (
            <button
              type="button"
              onClick={selectedIds.length > 0 ? handleDeleteSelectedUrls : handleDeleteAllUrls}
              className={`absolute right-0 -top-10 px-4 py-2 rounded text-white ${selectedIds.length > 0 ? 'bg-red-500' : 'bg-red-500'
                }`}
            >
              {selectedIds.length > 0 ? 'Delete' : 'Delete All'}
            </button>
          )}

          <table className="w-full border mt-4 border-blue-200">
            <thead className="bg-blue-100">
              <tr>
                <th className="p-2 border">Select</th>
                <th className="p-2 border">Long URL(s)</th>
                <th className="p-2 border">Short URL(s)</th>
                <th className="p-2 border">Actions</th>
              </tr>
            </thead>
            <tbody>
              {urls.map((url) => (
                <tr key={url.id} className="text-center border-t">
                  <td className="border p-2">
                    <input
                      type="checkbox"
                      checked={selectedIds.includes(url.id)}
                      onChange={() => handleSelectRow(url.id)}
                    />
                  </td>
                  <td className="border p-2">{url.longurl}</td>
                  <td className="border p-2">{url.shorturl}</td>
                  <td className="border p-2 flex justify-center space-x-2">
                    <button
                      onClick={() => handleDelete(url.id)}
                      className="bg-red-500 text-white px-3 py-1 rounded flex items-center space-x-1"
                    >
                      <Trash2 size={14} /> <span>Delete</span>
                    </button>
                    <button
                      onClick={() => handleEdit(url)}
                      className="bg-[brown] text-white px-3 py-1 rounded flex items-center space-x-1"
                    >
                      <Edit size={14} /> <span>Edit</span>
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </main>

      {editUrl && (
          <div className="fixed inset-0 bg-black/40 flex items-center justify-center z-50">
            <div className="bg-white rounded-lg shadow-lg p-6 w-96 relative">
              <button
                className="absolute top-2 right-2 text-gray-600"
                onClick={() => setEditUrl(null)}
              >
                <X />
              </button>
              <h3 className="text-lg font-semibold mb-4">Edit URL</h3>
              <input
                value={editUrl.shorturl}
                onChange={(e) => setEditUrl({ ...editUrl, shorturl: e.target.value })}
                className="border border-blue-300 rounded px-3 py-2 w-full mb-4"
              />
              <button
                onClick={handleSaveEdit}
                className="bg-blue-500 text-white px-4 py-2 rounded w-full"
              >
                Save Changes
              </button>
            </div>
          </div>
        )
      }
    </div>
  )
}