import RouteSearch from "../components/RouteSearch";

export default function SearchPage() {
  return (
    <div className="mx-auto max-w-lg p-4 md:max-w-2xl">
      <h1 className="mb-4 text-xl font-bold">搜尋路線</h1>
      <RouteSearch />
    </div>
  );
}
