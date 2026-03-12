# Lua सपोर्ट

[English](../../languages/lua.md) | [한국어](../../ko/languages/lua.md) | [日本語](../../ja/languages/lua.md) | [हिन्दी](lua.md) | [Deutsch](../../de/languages/lua.md)

## समर्थित एक्सटेंशन

- `.lua`

## एक्सट्रैक्शन लक्ष्य

| तत्व | Kind | उदाहरण |
|------|------|--------|
| ग्लोबल फ़ंक्शन | `function` | `function globalFunc(a, b)` |
| लोकल फ़ंक्शन | `local_function` | `local function helper()` |
| मॉड्यूल फ़ंक्शन | `module_function` | `function M.create(name)` |
| मेथड | `method` | `function M:init(config)` |
| टेबल असाइनमेंट | `variable` | `local M = {}` |
| फ़ंक्शन असाइनमेंट | `variable` | `local foo = function() end` |
| LuaDoc कमेंट | `doc` | `--- विवरण` |
| लाइन कमेंट | `doc` | `-- विवरण` |
| require() | `import` | `local json = require("json")` |

## उदाहरण

### इनपुट

```lua
local M = {}

--- नाम से किसी व्यक्ति का अभिवादन करता है।
-- @param name string व्यक्ति का नाम
function M.greet(name)
    print("Hello, " .. name)
end

function M:init(config)
    self.config = config
end

local function helper()
    return 42
end

function globalFunc(a, b)
    return a + b
end

local json = require("json")
```

### आउटपुट (XML)

```xml
<file path="example.lua" language="lua">
  <variable kind="variable" line="1">
    <name>M</name>
    <text>local M = {}</text>
  </variable>
  <function kind="module_function" line="4">
    <name>greet</name>
    <text>function M.greet(name)</text>
  </function>
  <function kind="method" line="9">
    <name>init</name>
    <text>function M:init(config)</text>
  </function>
  <function kind="local_function" line="13">
    <name>helper</name>
    <text>local function helper()</text>
  </function>
  <function kind="function" line="17">
    <name>globalFunc</name>
    <text>function globalFunc(a, b)</text>
  </function>
</file>
```

## नोट्स

### दृश्यता (Visibility)

- सभी डिक्लेरेशन एक्सट्रैक्ट किए जाते हैं (Lua में एक्सेस मॉडिफ़ायर नहीं हैं)
- `local` फ़ंक्शन/वेरिएबल भी ग्लोबल के साथ शामिल किए जाते हैं

### फ़ंक्शन प्रकार

- `function`: ग्लोबल फ़ंक्शन डिक्लेरेशन (`function foo()`)
- `local_function`: लोकल फ़ंक्शन डिक्लेरेशन (`local function foo()`)
- `module_function`: डॉट-इंडेक्स्ड फ़ंक्शन (`function M.foo()`)
- `method`: कोलन-इंडेक्स्ड मेथड (`function M:foo()`)

### बॉडी रिमूवल

`--include-body` फ़्लैग का उपयोग न करने पर:

- फ़ंक्शन/मेथड: पैरामीटर लिस्ट की क्लोजिंग ब्रैकेट `)` के बाद बॉडी हटा दी जाती है
- टेबल असाइनमेंट (`local M = {}`): यथावत रखा जाता है
- फ़ंक्शन असाइनमेंट (`local foo = function() end`): यथावत रखा जाता है

### इम्पोर्ट एक्सट्रैक्शन

- `require()` कॉल `--include-imports` फ़्लैग से एक्सट्रैक्ट किए जा सकते हैं
- `--include-private` का उपयोग करके गैर-निर्यातित/निजी सिंबल शामिल करें
- फ़ॉर्मेट: `local json = require("json")` (पूरा स्टेटमेंट संरक्षित रहता है)

### डॉक्यूमेंटेशन कमेंट्स

- LuaDoc कमेंट (`---`) और सामान्य लाइन कमेंट (`--`) दोनों एक्सट्रैक्ट किए जाते हैं
- ब्लॉक कमेंट (`--[[ ... ]]`) भी समर्थित हैं
