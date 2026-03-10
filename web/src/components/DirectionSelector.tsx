interface Props {
  direction: number;
  onChange: (d: number) => void;
}

export default function DirectionSelector({ direction, onChange }: Props) {
  const options = [
    { value: 0, label: "去程" },
    { value: 1, label: "回程" },
  ];

  return (
    <div className="flex gap-2" role="tablist">
      {options.map((opt) => (
        <button
          key={opt.value}
          type="button"
          role="tab"
          aria-selected={direction === opt.value}
          className={`flex-1 rounded-lg px-4 py-2 font-medium ${
            direction === opt.value
              ? "bg-blue-600 text-white"
              : "bg-gray-100 text-gray-600 hover:bg-gray-200"
          }`}
          onClick={() => onChange(opt.value)}
        >
          {opt.label}
        </button>
      ))}
    </div>
  );
}
