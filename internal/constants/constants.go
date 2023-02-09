package constants

import "time"

// Environment variable names
const ENV_CENSYS_API_ID = "CENSYS_API_ID"
const ENV_CENSYS_API_SECRET = "CENSYS_API_SECRET"

// Censys client defaults
const CENSYS_BASE_URL = "https://search.censys.io"
const CENSYS_REQ_SLEEP = 200 * time.Millisecond

const USE_LOCAL_IPLOOKUP = false
const SAVE_TO_FILE_INTERVAL = 5
