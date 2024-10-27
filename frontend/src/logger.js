// logger.js
import log from "loglevel";
import prefix from "loglevel-plugin-prefix";

// Apply the prefix plugin to loglevel
prefix.reg(log);
log.setLevel("info"); // Set the default log level for all components

function createComponentLogger(name, level = "info") {
  const componentLogger = log.getLogger(name);
  prefix.apply(componentLogger, {
    format(level, name, timestamp) {
      return `[${name}] [${level.toUpperCase()}] ${timestamp}`;
    },
    timestampFormatter(date) {
      return date.toISOString();
    },
  });

  componentLogger.setLevel(level); // Set the specific log level for this logger
  return componentLogger;
}

export { createComponentLogger };
