import { Component, type ReactNode } from "react";

interface Props {
  children: ReactNode;
}

interface State {
  error: Error | null;
}

export default class ErrorBoundary extends Component<Props, State> {
  state: State = { error: null };

  static getDerivedStateFromError(error: Error): State {
    return { error };
  }

  componentDidCatch(error: Error, info: React.ErrorInfo) {
    console.error("[ErrorBoundary]", error, info.componentStack);
  }

  render() {
    if (this.state.error) {
      return (
        <div className="mx-auto max-w-lg p-6 text-center">
          <h1 className="mb-2 text-xl font-bold text-red-600">發生錯誤</h1>
          <p className="mb-4 text-sm text-gray-600 dark:text-gray-400">
            {this.state.error.message}
          </p>
          <button
            type="button"
            className="rounded bg-blue-600 px-4 py-2 text-white hover:bg-blue-700"
            onClick={() => this.setState({ error: null })}
          >
            重試
          </button>
        </div>
      );
    }
    return this.props.children;
  }
}
