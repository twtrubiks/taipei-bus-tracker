import { Link } from "react-router-dom";

export default function HomePage() {
  return (
    <div className="mx-auto max-w-lg p-4">
      <h1 className="mb-6 text-2xl font-bold">公車到站查詢</h1>
      <Link
        to="/search"
        className="flex items-center gap-2 rounded-lg bg-blue-600 px-4 py-3 text-white shadow hover:bg-blue-700"
      >
        <span className="text-xl">&#128269;</span>
        搜尋路線
      </Link>
    </div>
  );
}
