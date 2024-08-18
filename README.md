# hc

Habbo Codec

## Usage

## VL64

### Encode

```sh
$ hc -vl64 {1..3}
1: I
2: J
3: K

# with -c: compact
$ hc -vl64e -c {0..30..6}
HRAPCRDPFRG
```

### Decode

```sh
$ hc -vl64 HRAPCRDPFRG
H: 0
RA: 6
PC: 12
RD: 18
PF: 24
RG: 30
```

## B64

### Encode

```sh
$ hc -b64e 53
@u
```

### Decode

```sh
$ hc -b64 @u
53
```