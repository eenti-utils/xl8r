# xl8r

The xl8r package is a library written in Go that facilitates development of "spoke and hub" translators.

**Spoke and Hub:**
```mermaid
flowchart TB;
  X((Point<br>Px))<-->H(((Hub)));
  A((Point<br>P1))<-->H;
  B((Point<br>P2))<-->H;
  C((Point<br>P3))<-->H;
  N((Point<br>Pn))<-->H;

  style H fill:#FFD9B8,stroke:black,color:black
  style X fill:#F0F0F0,stroke:#0000FF,color:#0000FF
  style A fill:#F0F0F0,stroke:#FF9900,color:#FF9900
  style B fill:#F0F0F0,stroke:#cc00ff,color:#cc00ff
  style C fill:#F0F0F0,stroke:#663300,color:#663300
  style N fill:#F0F0F0,stroke:#00F000,color:#00F000

```
<details>

<summary>More on xl8r "spoke and hub" paradigm</summary>

In this paradigm, each point represents a different Origin and/or Destination for data translations.

For example, some content may be translated from `Point P1` to `Point P3`, where:
- `Point P1` is called "english"
- `Point P3` is called "spanish"
- the content to be translated is the value `string` "four"

In this example (and the `xl8r` package), `Point P1` is considered as the _Origin_ and `Point P3` as the _Destination_.

The _Hub_ represents a commonality between _all points_ in the system.
- the hub data, in this example, is the value `int` 4

The _Spoke_ represents the path to and from `Point` and `Hub`.
- from `Point` (_Origin_) to `Hub`, "content data" is converted to "hub data"  (ie. _Encoded_)
- from `Hub` to `Point` (_Destination_), "hub data" is converted to "content data"  (ie. _Decoded_)

Summarizing the "english" to "spanish" translation, in _spoke and hub_ terms:
- from "english": (`Point P1`) convert value `string` "four" to value `int` 4 (`Hub`)
- to "spanish": (`Hub`) convert value `int` 4 to value `string` "cuatro" (`Point P3`)

</details>

```mermaid
graph LR
  Dx(Content Data<br>Px) --> X((Point<br>Px))
  X -- Encode --> H(((Hub Data)))
  H -- Decode --> N((Point<br>Pn))
  N --> Dn(Content Data<br>Pn)

  style H fill:#FFD9B8,stroke:black,color:black
  style Dx fill:#F0F0F0,stroke:#0000FF,color:#0000FF
  style X fill:#F0F0F0,stroke:#0000FF,color:#0000FF
  style N fill:#F0F0F0,stroke:#00F000,color:#00F000
  style Dn fill:#F0F0F0,stroke:#00F000,color:#00F000
```

In 

```mermaid
graph LR
  A((Point)) -- Encode --> H(((Hub Data)))
  H -- Decode --> A((Point))

  style H fill:#FFD9B8,stroke:black,color:black
```

