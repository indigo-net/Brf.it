# Ruby सपोर्ट

[English](../../languages/ruby.md) | [한국어](../../ko/languages/ruby.md) | [日本語](../../ja/languages/ruby.md) | [हिन्दी](ruby.md) | [Deutsch](../../de/languages/ruby.md)

## समर्थित एक्सटेंशन

- `.rb`

## एक्सट्रैक्शन लक्ष्य

| तत्व | Kind | XML Tag | उदाहरण |
|------|------|---------|--------|
| मेथड | `method` | `<function>` | `def greet(name)` |
| क्लास मेथड | `method` | `<function>` | `def self.create(attrs)` |
| क्लास | `class` | `<type>` | `class User < ActiveRecord::Base` |
| मॉड्यूल | `namespace` | `<type>` | `module Authentication` |
| कॉन्स्टेंट (top-level) | `variable` | `<variable>` | `MAX_RETRIES = 3` |
| YARD कमेंट | `doc` | | `# विवरण` |
| require | | `<imports>` | `require "json"` |
| require_relative | | `<imports>` | `require_relative "helpers"` |

## उदाहरण

### इनपुट

```ruby
require "json"
require_relative "helpers"

MAX_RETRIES = 3

# सिस्टम में एक उपयोगकर्ता को दर्शाता है।
class User
  # एट्रिब्यूट्स से एक नया उपयोगकर्ता बनाता है।
  def self.create(attrs)
    new(attrs).save
  end

  # उपयोगकर्ता को इनिशियलाइज़ करता है।
  def initialize(name)
    @name = name
  end

  # किसी अन्य व्यक्ति का अभिवादन करता है।
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
  <type>class User</type>
  <function>def self.create(attrs)</function>
  <function>def initialize(name)</function>
  <function>def greet(other)</function>
  <variable>MAX_RETRIES = 3</variable>
  <type>module Authentication</type>
  <function>def authenticate(password)</function>
</file>
```

## नोट्स

### दृश्यता (Visibility)

- विजिबिलिटी (`public`, `protected`, `private`) की परवाह किए बिना सभी मेथड एक्सट्रैक्ट किए जाते हैं
- इंस्टेंस मेथड (`def foo`) और क्लास मेथड (`def self.foo`) दोनों कैप्चर किए जाते हैं

### मेथड प्रकार

- इंस्टेंस मेथड (`def foo`) और क्लास मेथड (`def self.foo`) दोनों kind `method` का उपयोग करते हैं

### बॉडी रिमूवल

`--include-body` फ़्लैग का उपयोग न करने पर:

- मेथड: पैरामीटर लिस्ट की क्लोजिंग ब्रैकेट `)` के बाद बॉडी हटा दी जाती है (पैरामीटर न होने पर मेथड नाम के बाद)
- क्लास/मॉड्यूल: केवल डिक्लेरेशन लाइन संरक्षित रहती है
- कॉन्स्टेंट: यथावत रखे जाते हैं

### इम्पोर्ट एक्सट्रैक्शन

- `require` और `require_relative` स्टेटमेंट `--include-imports` फ़्लैग से एक्सट्रैक्ट किए जा सकते हैं
- `--include-private` का उपयोग करके गैर-निर्यातित/निजी सिंबल शामिल करें
- फ़ॉर्मेट: `require "json"` / `require_relative "helpers"` (पूरा स्टेटमेंट संरक्षित रहता है)

### डॉक्यूमेंटेशन कमेंट्स

- YARD स्टाइल कमेंट (`#`) जो मेथड/क्लास के ठीक ऊपर हों, एक्सट्रैक्ट किए जाते हैं
- मल्टी-लाइन `#` कमेंट भी समर्थित हैं
- `=begin`...`=end` ब्लॉक कमेंट भी पहचाने जाते हैं
