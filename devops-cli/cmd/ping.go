package cmd

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

// pingCmd sends an HTTP GET request to the given URL and reports status & latency.
// It is useful for quick uptime checks or simple API responsiveness tests.
var pingCmd = &cobra.Command{
	Use:   "ping [url]",
	Short: "Ping a URL via HTTP GET and measure response time",
	Long: `Send an HTTP GET request to a given URL and print:
- HTTP status line (e.g., 200 OK)
- Total round-trip duration (wall-clock latency)
This is handy for lightweight health checks.`,
	Args: cobra.ExactArgs(1), // Require exactly one argument: the target URL.
	RunE: func(cmd *cobra.Command, args []string) error {
		url := args[0]

		// Read flags (timeout, follow-redirects) with sane defaults.
		timeout, _ := cmd.Flags().GetDuration("timeout")
		followRedirects, _ := cmd.Flags().GetBool("follow-redirects")

		// Create a cancellable context to enforce request timeout.
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		// Build an HTTP client; optionally prevent following redirects.
		client := &http.Client{
			Timeout: timeout,
			// Control redirect policy: return on first redirect if disabled.
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				if followRedirects {
					return nil // allow
				}
				// Returning an error stops following redirects (documented behavior).
				return http.ErrUseLastResponse
			},
		}

		// Create the request with context to ensure proper timeout cancellation.
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return fmt.Errorf("build request: %w", err)
		}

		// Measure wall-clock duration from before Do() until response/error.
		start := time.Now()
		resp, err := client.Do(req)
		elapsed := time.Since(start)

		// Network/timeout/DNS errors are returned here.
		if err != nil {
			return fmt.Errorf("request error: %w (elapsed=%s)", err, elapsed)
		}
		// Ensure response body is closed to free resources.
		defer resp.Body.Close()

		// Print human-friendly status line and latency.
		fmt.Printf("✅ Status: %s\n", resp.Status)
		fmt.Printf("⏱️ Response Time: %v\n", elapsed)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(pingCmd)

	// --timeout: maximum total time for the HTTP request (connect + TLS + read).
	pingCmd.Flags().Duration("timeout", 10*time.Second, "Overall request timeout")

	// --follow-redirects: follow 3xx responses (e.g., 301/302). Default false to
	// expose redirect behavior explicitly for diagnostics.
	pingCmd.Flags().Bool("follow-redirects", false, "Follow HTTP redirects")
}
