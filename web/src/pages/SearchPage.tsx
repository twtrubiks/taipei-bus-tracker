import { Link } from "react-router-dom";
import RouteSearch from "../components/RouteSearch";

export default function SearchPage() {
  return (
    <div className="mx-auto max-w-lg p-4">
      <div className="mb-4 flex items-center gap-2">
        <Link to="/" className="text-blue-600 hover:text-blue-800 dark:text-blue-400 dark:hover:text-blue-300">
          &larr; 返回
        </Link>
        <h1 className="text-xl font-bold">搜尋路線</h1>
      </div>
      <RouteSearch />
    </div>
  );
}
