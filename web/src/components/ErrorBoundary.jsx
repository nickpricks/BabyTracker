import React from "react";

export default class ErrorBoundary extends React.Component {
  constructor(props) {
    super(props);
    this.state = { hasError: false, error: null };
  }

  static getDerivedStateFromError(error) {
    return { hasError: true, error };
  }

  render() {
    if (this.state.hasError) {
      return (
        <div className="card text-center py-12 max-w-md mx-auto mt-10">
          <p className="text-4xl mb-4">😵</p>
          <h2 className="font-display text-xl font-bold text-fg-heading mb-2">
            Something went wrong
          </h2>
          <p className="text-sm text-fg-muted mb-6">
            The app encountered an unexpected error. Try refreshing the page.
          </p>
          <button
            onClick={() => window.location.reload()}
            className="btn-primary"
          >
            Refresh
          </button>
        </div>
      );
    }
    return this.props.children;
  }
}
