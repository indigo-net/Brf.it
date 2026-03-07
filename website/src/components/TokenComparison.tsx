import React, {useState} from 'react';
import Translate, {translate} from '@docusaurus/Translate';

interface CodeExample {
  language: string;
  label: string;
  before: string;
  after: string;
  beforeTokens: number;
  afterTokens: number;
  beforeLines: number;
  afterLines: number;
}

const examples: CodeExample[] = [
  {
    language: 'typescript',
    label: 'TypeScript',
    before: `import { useState, useEffect } from 'react';
import axios from 'axios';

interface User {
  id: number;
  name: string;
  email: string;
  role: 'admin' | 'user' | 'guest';
}

interface UserService {
  fetchUser(id: number): Promise<User>;
  updateUser(user: User): Promise<User>;
  deleteUser(id: number): Promise<void>;
}

export class UserServiceImpl implements UserService {
  private baseUrl = '/api/users';

  async fetchUser(id: number): Promise<User> {
    const response = await axios.get(\`\${this.baseUrl}/\${id}\`);
    return response.data;
  }

  async updateUser(user: User): Promise<User> {
    const response = await axios.put(\`\${this.baseUrl}/\${user.id}\`, user);
    return response.data;
  }

  async deleteUser(id: number): Promise<void> {
    await axios.delete(\`\${this.baseUrl}/\${id}\`);
  }
}

export function useUser(userId: number) {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);

  useEffect(() => {
    const service = new UserServiceImpl();
    service.fetchUser(userId)
      .then(setUser)
      .catch(setError)
      .finally(() => setLoading(false));
  }, [userId]);

  return { user, loading, error };
}`,
    after: `// UserServiceImpl.ts
import { useState, useEffect } from 'react';
import axios from 'axios';

interface User { id: number; name: string; email: string; role: 'admin' | 'user' | 'guest' }
interface UserService { fetchUser(id: number): Promise<User>; updateUser(user: User): Promise<User>; deleteUser(id: number): Promise<void> }
class UserServiceImpl implements UserService { /* ... */ }
function useUser(userId: number): { user: User | null; loading: boolean; error: Error | null }`,
    beforeTokens: 420,
    afterTokens: 95,
    beforeLines: 52,
    afterLines: 8,
  },
  {
    language: 'python',
    label: 'Python',
    before: `from typing import List, Optional, Dict, Any
from dataclasses import dataclass
from datetime import datetime
import asyncio

@dataclass
class User:
    id: int
    name: str
    email: str
    created_at: datetime
    metadata: Dict[str, Any]

class UserRepository:
    def __init__(self, db_connection):
        self.db = db_connection

    async def get_user(self, user_id: int) -> Optional[User]:
        query = "SELECT * FROM users WHERE id = $1"
        row = await self.db.fetchrow(query, user_id)
        if row:
            return User(
                id=row['id'],
                name=row['name'],
                email=row['email'],
                created_at=row['created_at'],
                metadata=row['metadata']
            )
        return None

    async def list_users(self, limit: int = 100, offset: int = 0) -> List[User]:
        query = "SELECT * FROM users LIMIT $1 OFFSET $2"
        rows = await self.db.fetch(query, limit, offset)
        return [User(**row) for row in rows]

    async def create_user(self, name: str, email: str) -> User:
        query = """
            INSERT INTO users (name, email, created_at, metadata)
            VALUES ($1, $2, NOW(), '{}')
            RETURNING *
        """
        row = await self.db.fetchrow(query, name, email)
        return User(**row)`,
    after: `# user_repository.py
from typing import List, Optional, Dict, Any
from dataclasses import dataclass
from datetime import datetime

@dataclass
class User: id: int; name: str; email: str; created_at: datetime; metadata: Dict[str, Any]

class UserRepository:
    def __init__(self, db_connection): ...
    async def get_user(self, user_id: int) -> Optional[User]: ...
    async def list_users(self, limit: int = 100, offset: int = 0) -> List[User]: ...
    async def create_user(self, name: str, email: str) -> User: ...`,
    beforeTokens: 380,
    afterTokens: 78,
    beforeLines: 44,
    afterLines: 12,
  },
  {
    language: 'go',
    label: 'Go',
    before: `package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrUnauthorized = errors.New("unauthorized")
)

type User struct {
	ID        string    \`json:"id"\`
	Name      string    \`json:"name"\`
	Email     string    \`json:"email"\`
	CreatedAt time.Time \`json:"created_at"\`
	UpdatedAt time.Time \`json:"updated_at"\`
}

type UserService struct {
	db    Database
	cache *redis.Client
	http  *http.Client
}

func NewUserService(db Database, cache *redis.Client) *UserService {
	return &UserService{
		db:    db,
		cache: cache,
		http:  &http.Client{Timeout: 30 * time.Second},
	}
}

func (s *UserService) GetUser(ctx context.Context, id string) (*User, error) {
	cacheKey := fmt.Sprintf("user:%s", id)
	cached, err := s.cache.Get(ctx, cacheKey).Result()
	if err == nil {
		var user User
		if err := json.Unmarshal([]byte(cached), &user); err == nil {
			return &user, nil
		}
	}

	user, err := s.db.FindUser(ctx, id)
	if err != nil {
		return nil, ErrUserNotFound
	}

	data, _ := json.Marshal(user)
	s.cache.Set(ctx, cacheKey, data, 5*time.Minute)

	return user, nil
}`,
    after: `// service/user.go
package service

import (
	"context"
	"encoding/json"
	"time"
	"github.com/redis/go-redis/v9"
)

var ErrUserNotFound = errors.New("user not found")

type User struct { ID string; Name string; Email string; CreatedAt time.Time; UpdatedAt time.Time }
type UserService struct { db Database; cache *redis.Client; http *http.Client }
func NewUserService(db Database, cache *redis.Client) *UserService
func (s *UserService) GetUser(ctx context.Context, id string) (*User, error)`,
    beforeTokens: 450,
    afterTokens: 85,
    beforeLines: 65,
    afterLines: 13,
  },
];

export default function TokenComparison(): JSX.Element {
  const [selectedLanguage, setSelectedLanguage] = useState(examples[0].language);
  const currentExample = examples.find(e => e.language === selectedLanguage)!;

  const tokenReduction = Math.round(
    ((currentExample.beforeTokens - currentExample.afterTokens) / currentExample.beforeTokens) * 100
  );
  const lineReduction = Math.round(
    ((currentExample.beforeLines - currentExample.afterLines) / currentExample.beforeLines) * 100
  );

  return (
    <section className="section">
      <div className="container">
        <h2 className="section-title">
          <Translate id="comparison.title">Token Savings in Action</Translate>
        </h2>
        <p className="section-subtitle">
          <Translate id="comparison.subtitle">
            See how Brf.it reduces token usage while preserving essential context
          </Translate>
        </p>

        {/* Language Tabs */}
        <div className="tabs-container" style={{ display: 'flex', justifyContent: 'center', gap: '0.5rem', marginBottom: '2rem' }}>
          {examples.map(example => (
            <button
              key={example.language}
              onClick={() => setSelectedLanguage(example.language)}
              className={`tab-button ${selectedLanguage === example.language ? 'active' : ''}`}
              style={{
                padding: '0.5rem 1.5rem',
                border: selectedLanguage === example.language ? '2px solid var(--ifm-color-primary)' : '1px solid var(--ifm-color-emphasis-400)',
                borderRadius: '8px',
                background: selectedLanguage === example.language ? 'var(--ifm-color-primary)' : 'transparent',
                color: selectedLanguage === example.language ? 'white' : 'var(--ifm-color-emphasis-700)',
                cursor: 'pointer',
                fontWeight: 500,
              }}
            >
              {example.label}
            </button>
          ))}
        </div>

        {/* Code Comparison */}
        <div className="comparison-grid">
          <div className="code-panel">
            <div className="code-panel-header">
              <span className="code-panel-title">Before (Full Code)</span>
              <span className="code-panel-badge before">
                {currentExample.beforeTokens.toLocaleString()} tokens
              </span>
            </div>
            <pre className="code-panel-content">
              <code>{currentExample.before}</code>
            </pre>
          </div>

          <div className="code-panel">
            <div className="code-panel-header">
              <span className="code-panel-title">After (Brf.it)</span>
              <span className="code-panel-badge after">
                {currentExample.afterTokens.toLocaleString()} tokens
              </span>
            </div>
            <pre className="code-panel-content">
              <code>{currentExample.after}</code>
            </pre>
          </div>
        </div>

        {/* Stats */}
        <div className="stats-card">
          <div className="stats-grid">
            <div className="stat-item">
              <div className="stat-value">{tokenReduction}%</div>
              <div className="stat-label">
                <Translate id="comparison.tokens">Token Reduction</Translate>
              </div>
              <div className="stat-change positive">
                -{(currentExample.beforeTokens - currentExample.afterTokens).toLocaleString()} tokens
              </div>
            </div>
            <div className="stat-item">
              <div className="stat-value">{lineReduction}%</div>
              <div className="stat-label">
                <Translate id="comparison.lines">Fewer Lines</Translate>
              </div>
              <div className="stat-change positive">
                -{currentExample.beforeLines - currentExample.afterLines} lines
              </div>
            </div>
            <div className="stat-item">
              <div className="stat-value">~${((currentExample.beforeTokens - currentExample.afterTokens) * 0.000003).toFixed(2)}</div>
              <div className="stat-label">
                <Translate id="comparison.savings">Savings per Request</Translate>
              </div>
              <div className="stat-change positive">
                @ $3/1M tokens
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}
