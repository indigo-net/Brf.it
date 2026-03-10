use std::collections::HashMap;
use std::fmt;

/// A key-value store with expiration support.
pub struct Cache<V> {
    data: HashMap<String, V>,
    capacity: usize,
}

/// Errors that can occur during cache operations.
#[derive(Debug)]
pub enum CacheError {
    KeyNotFound(String),
    CapacityExceeded,
}

impl fmt::Display for CacheError {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        match self {
            CacheError::KeyNotFound(key) => write!(f, "key not found: {}", key),
            CacheError::CapacityExceeded => write!(f, "cache capacity exceeded"),
        }
    }
}

/// Trait for items that can be cached.
pub trait Cacheable: Clone + fmt::Debug {
    /// Returns the size in bytes of this item.
    fn size(&self) -> usize;
}

impl<V: Clone> Cache<V> {
    /// Create a new cache with the given capacity.
    pub fn new(capacity: usize) -> Self {
        Cache {
            data: HashMap::new(),
            capacity,
        }
    }

    /// Get a value by key. Returns None if not found.
    pub fn get(&self, key: &str) -> Option<&V> {
        self.data.get(key)
    }

    /// Insert a key-value pair. Returns error if capacity exceeded.
    pub fn insert(&mut self, key: String, value: V) -> Result<(), CacheError> {
        if self.data.len() >= self.capacity && !self.data.contains_key(&key) {
            return Err(CacheError::CapacityExceeded);
        }
        self.data.insert(key, value);
        Ok(())
    }

    /// Remove a key and return its value.
    pub fn remove(&mut self, key: &str) -> Result<V, CacheError> {
        self.data
            .remove(key)
            .ok_or_else(|| CacheError::KeyNotFound(key.to_string()))
    }

    /// Returns the number of items in the cache.
    pub fn len(&self) -> usize {
        self.data.len()
    }

    /// Returns true if the cache is empty.
    pub fn is_empty(&self) -> bool {
        self.data.is_empty()
    }
}
