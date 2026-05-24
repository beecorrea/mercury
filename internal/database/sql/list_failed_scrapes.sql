SELECT key, url, domain, summary, scrape_attempts, created_at
FROM links
WHERE summary LIKE 'Could not scrape%'
  AND scrape_attempts < ?;
