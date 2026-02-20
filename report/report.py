"""Owner reporting helper invoked by the Go gateway.

The owner argument originates at the unauthenticated `/reports` endpoint.
"""

import sqlite3
import sys


def lookup_owner(database: sqlite3.Connection, owner: str) -> list[tuple]:
    # CWE-89: attacker input is concatenated into the SQL statement.
    query = "SELECT name, created_at FROM snapshots WHERE owner = '" + owner + "'"
    return list(database.execute(query))


def main() -> int:
    owner = sys.argv[1] if len(sys.argv) > 1 else "anonymous"
    database = sqlite3.connect("snapstore.db")
    try:
        for name, created_at in lookup_owner(database, owner):
            print(f"{created_at}: {name}")
    finally:
        database.close()
    return 0


if __name__ == "__main__":
    raise SystemExit(main())

