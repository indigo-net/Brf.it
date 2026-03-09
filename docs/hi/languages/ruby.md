# Ruby सपोर्ट

[English](../../languages/ruby.md) | [한국어](../../ko/languages/ruby.md) | [日本語](../../ja/languages/ruby.md) | [हिन्दी](ruby.md) | [Deutsch](../../de/languages/ruby.md)

## समर्थित एक्सटेंशन

- `.rb`

## एक्सट्रैक्शन लक्ष्य

| तत्व | Kind | उदाहरण |
|------|------|--------|
| मेथड | `method` | `def greet(name)` |
| क्लास मेथड | `class_method` | `def self.create(attrs)` |
| क्लास | `class` | `class User < ActiveRecord::Base` |
| मॉड्यूल | `module` | `module Authentication` |
| कॉन्स्टेंट | `variable` | `MAX_RETRIES = 3` |
| YARD कमेंट | `doc` | `# विवरण` |
| require | `import` | `require "json"` |
| require_relative | `import` | `require_relative "helpers"` |

## उदाहरण

### इनपुट

```ruby
require "json"
require_relative "helpers"

# सिस्टम में एक उपयोगकर्ता को दर्शाता है।
class User
  MAX_RETRIES = 3

  # एट्रिब्यूट्स से एक नया उपयोगकर्ता बनाता है।
  # @param attrs [Hash] उपयोगकर्ता एट्रिब्यूट्स
  def self.create(attrs)
    new(attrs).save
  end

  # उपयोगकर्ता को इनिशियलाइज़ करता है।
  # @param name [String] उपयोगकर्ता का नाम
  def initialize(name)
    @name = name
  end

  # किसी अन्य व्यक्ति का अभिवादन करता है।
  # @param other [String] दूसरे व्यक्ति का नाम
  # @return [String] अभिवादन संदेश
  def greet(other)
    "Hello, #{other}! I'm #{@name}."
  end
end

module Authentication
  def authenticate(password)
    password == @secret
  end
end
```

### आउटपुट (XML)

```xml
<file path="example.rb" language="ruby">
  <class kind="class" line="5">
    <name>User</name>
    <text>class User</text>
  </class>
  <variable kind="variable" line="6">
    <name>MAX_RETRIES</name>
    <text>MAX_RETRIES = 3</text>
  </variable>
  <function kind="class_method" line="10">
    <name>create</name>
    <text>def self.create(attrs)</text>
  </function>
  <function kind="method" line="15">
    <name>initialize</name>
    <text>def initialize(name)</text>
  </function>
  <function kind="method" line="21">
    <name>greet</name>
    <text>def greet(other)</text>
  </function>
  <module kind="module" line="27">
    <name>Authentication</name>
    <text>module Authentication</text>
  </module>
  <function kind="method" line="28">
    <name>authenticate</name>
    <text>def authenticate(password)</text>
  </function>
</file>
```

## नोट्स

### दृश्यता (Visibility)

- विजिबिलिटी (`public`, `protected`, `private`) की परवाह किए बिना सभी मेथड एक्सट्रैक्ट किए जाते हैं
- इंस्टेंस मेथड (`def foo`) और क्लास मेथड (`def self.foo`) दोनों कैप्चर किए जाते हैं

### मेथड प्रकार

- `method`: इंस्टेंस मेथड डिक्लेरेशन (`def foo`)
- `class_method`: क्लास-लेवल मेथड डिक्लेरेशन (`def self.foo`)

### बॉडी रिमूवल

`--include-body` फ़्लैग का उपयोग न करने पर:

- मेथड: पैरामीटर लिस्ट की क्लोजिंग ब्रैकेट `)` के बाद बॉडी हटा दी जाती है (पैरामीटर न होने पर मेथड नाम के बाद)
- क्लास/मॉड्यूल: केवल डिक्लेरेशन लाइन संरक्षित रहती है
- कॉन्स्टेंट: यथावत रखे जाते हैं

### इम्पोर्ट एक्सट्रैक्शन

- `require` और `require_relative` स्टेटमेंट `--include-imports` फ़्लैग से एक्सट्रैक्ट किए जा सकते हैं
- फ़ॉर्मेट: `require "json"` / `require_relative "helpers"` (पूरा स्टेटमेंट संरक्षित रहता है)

### डॉक्यूमेंटेशन कमेंट्स

- YARD स्टाइल कमेंट (`#`) जो मेथड/क्लास के ठीक ऊपर हों, एक्सट्रैक्ट किए जाते हैं
- मल्टी-लाइन `#` कमेंट भी समर्थित हैं
- `=begin`...`=end` ब्लॉक कमेंट भी पहचाने जाते हैं
