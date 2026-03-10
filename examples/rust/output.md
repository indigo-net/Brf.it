# Code Summary: examples/rust

*brf.it dev*

---

## Files

### examples/rust/lib.rs

```rust
pub struct Cache<V>
pub enum CacheError
impl fmt::Display for CacheError
fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result
pub trait Cacheable: Clone + fmt::Debug
fn size(&self) -> usize;
impl<V: Clone> Cache<V>
pub fn new(capacity: usize) -> Self
pub fn get(&self, key: &str) -> Option<&V>
pub fn insert(&mut self, key: String, value: V) -> Result<(), CacheError>
pub fn remove(&mut self, key: &str) -> Result<V, CacheError>
pub fn len(&self) -> usize
pub fn is_empty(&self) -> bool
```

---

