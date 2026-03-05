# C#-Unterstützung

[English](../../languages/csharp.md) | [한국어](../../ko/languages/csharp.md) | [日本語](../../ja/languages/csharp.md) | [हिन्दी](../../hi/languages/csharp.md) | [Deutsch](csharp.md)

## Unterstützte Erweiterungen

- `.cs`

## Extraktionsziele

| Element | Kind | Beispiel |
|---------|------|----------|
| Klasse | `class` | `public class Calculator` |
| Struktur | `struct` | `public struct Point` |
| Schnittstelle | `interface` | `public interface IDrawable` |
| Aufzählung | `enum` | `public enum Color` |
| Record | `record` | `public record Person(string Name, int Age)` |
| Record-Struktur | `struct` | `public record struct Measurement(double Value)` |
| Delegat | `type` | `public delegate void Action<T>(T item)` |
| Methode | `method` | `public int Add(int a, int b)` |
| Konstruktor | `constructor` | `public Calculator()` |
| Destruktor | `destructor` | `~Calculator()` |
| Eigenschaft | `variable` | `public string Name { get; set; }` |
| Feld (static/const) | `variable` | `public const int MaxValue = 100` |
| Ereignis | `variable` | `public event EventHandler Changed` |
| Indexer | `method` | `public int this[int index]` |
| Operator | `function` | `public static operator +(...)` |
| Konvertierungsoperator | `function` | `public static implicit operator int(...)` |
| Namensraum | `namespace` | `namespace MyApp` |
| Aufzählungsmitglied | `variable` | `Red, Green, Blue` |
| Dokumentationskommentar | `doc` | `/// <summary>...</summary>` |

## Beispiel

### Eingabe

```csharp
using System;

namespace MyApp
{
    /// <summary>Calculator-Klasse für grundlegende Mathematik.</summary>
    public class Calculator
    {
        public const int MaxValue = 100;

        public Calculator() { }

        public int Add(int a, int b)
        {
            return a + b;
        }

        public string Name { get; set; }
    }

    public interface IService
    {
        void Execute();
    }

    public record Person(string Name, int Age);
}
```

### Ausgabe (XML)

```xml
<file path="Calculator.cs" language="csharp">
  <type kind="namespace" line="3">
    <name>MyApp</name>
    <text>namespace MyApp</text>
  </type>
  <type kind="class" line="6">
    <name>Calculator</name>
    <text>public class Calculator</text>
    <doc>Calculator-Klasse für grundlegende Mathematik.</doc>
  </type>
  <variable kind="variable" line="8">
    <name>MaxValue</name>
    <text>public const int MaxValue = 100;</text>
  </variable>
  <function kind="constructor" line="10">
    <name>Calculator</name>
    <text>public Calculator()</text>
  </function>
  <function kind="method" line="12">
    <name>Add</name>
    <text>public int Add(int a, int b)</text>
  </function>
  <variable kind="variable" line="17">
    <name>Name</name>
    <text>public string Name { get; set; }</text>
  </variable>
  <type kind="interface" line="20">
    <name>IService</name>
    <text>public interface IService</text>
  </type>
  <function kind="method" line="22">
    <name>Execute</name>
    <text>void Execute();</text>
  </function>
  <type kind="record" line="25">
    <name>Person</name>
    <text>public record Person(string Name, int Age);</text>
  </type>
</file>
```

## Hinweise

### Sichtbarkeit

- Alle Deklarationen werden unabhängig von Zugriffsmodifikatoren extrahiert
- Zugriffsmodifikatoren (`public`, `private`, `internal`, `protected`) bleiben in Signaturen erhalten

### Felder

- Nur `static`, `const` und `static readonly` Felder werden extrahiert
- Instanzfelder werden zur Rauschreduzierung ausgeschlossen

### Eigenschaften

- Auto-Eigenschaften (`{ get; set; }`) werden vollständig beibehalten
- Ausdruckskörper-Eigenschaften (`=> expr`) werden im Signaturmodus ohne Ausdruck zurückgegeben

### Records

- `record` und `record class` werden als Kind `record` klassifiziert
- `record struct` wird als Kind `struct` klassifiziert

### Operatoren

- Operator-Überladungen erhalten synthetisierte Namen wie `operator+`, `operator==`
- Konvertierungsoperatoren erhalten Namen wie `implicit operator int`, `explicit operator string`
- Indexer erhalten den synthetisierten Namen `this`

### Körperentfernung

Wenn das `--include-body` Flag nicht verwendet wird:

- Methoden/Konstruktoren/Destruktoren: Körper nach öffnender Klammer `{` entfernt
- Ausdruckskörper-Mitglieder: `=>` und Ausdruck entfernt
- Klassen/Strukturen/Schnittstellen/Aufzählungen/Records: Körper nach öffnender Klammer `{` entfernt
- Eigenschaften: Auto-Eigenschaften beibehalten, Ausdruckskörper-Eigenschaften entfernt
- Delegate: kein Körper, unverändert zurückgegeben

### Dokumentationskommentare

- Sowohl `///` XML-Dokumentationskommentare als auch `//` Zeilenkommentare werden extrahiert
- `/* */` Blockkommentare werden ebenfalls erfasst
