//Copyright 2019 Expedia, Inc.
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.

package logging

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"
)

// Setup initializes the global slog logger based on the provided configuration.
func Setup(level string, format string, timeFormat string) error {
	lvl, err := parseLevel(level)
	if err != nil {
		return err
	}

	replacer := buildTimeReplacer(timeFormat)

	opts := &slog.HandlerOptions{
		Level:       lvl,
		ReplaceAttr: replacer,
	}

	var handler slog.Handler
	switch strings.ToLower(format) {
	case "json":
		handler = slog.NewJSONHandler(os.Stderr, opts)
	case "text":
		handler = slog.NewTextHandler(os.Stderr, opts)
	default:
		return fmt.Errorf("unsupported log format %q, use text or json", format)
	}

	slog.SetDefault(slog.New(handler))
	return nil
}

func parseLevel(level string) (slog.Level, error) {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug, nil
	case "info":
		return slog.LevelInfo, nil
	case "warn", "warning":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	default:
		return slog.LevelInfo, fmt.Errorf("unsupported log level %q, use debug, info, warn, or error", level)
	}
}

func buildTimeReplacer(timeFormat string) func(groups []string, a slog.Attr) slog.Attr {
	var layout string
	switch strings.ToLower(timeFormat) {
	case "rfc3339":
		layout = time.RFC3339
	case "rfc3339nano":
		layout = time.RFC3339Nano
	case "epoch":
		layout = "epoch"
	default:
		// iso8601 or any unrecognized value defaults to RFC3339 (ISO 8601 compatible)
		layout = time.RFC3339
	}

	return func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey {
			t := a.Value.Time()
			if layout == "epoch" {
				a.Value = slog.Float64Value(float64(t.UnixMilli()) / 1000.0)
			} else {
				a.Value = slog.StringValue(t.Format(layout))
			}
		}
		return a
	}
}
