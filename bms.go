package main

/*

Mapping of BMS commands to output structs
*/

// bmsc> help
//   SHow [<>|Version|COnfig|MAP|CELLS|LTC|STATS|THermistors]
//         <>        - status
//         version   - firmware version
//         config    - configuration
//         map       - cell group map
//         cells     - cell summary
//         ltc       - measurement chips
//         stats     - cell statistics
//         thermistors - thermistor readings
//   SEt  [<>|ID|LVC|HVC|LVCC|HVCC|BVC|BVMIN|THMAX|MAP]
//         <>        - show config
//         id        - set bmsc ID (1..4)
//         lvc       - Low Voltage Cutoff
//         hvc       - High Voltage Cutoff
//         lvcc      - Low Voltage Cutoff Clear
//         hvcc      - High Voltage Cutoff Clear
//         bvc       - Charge balancing Voltage Cutoff
//         bvmin     - Auto Balancing Voltage Minimum
//         thmax     - Thermistor Max Temperature
//         map <ltc> <pack> <grp> - Map an LTC to a Pack/Cell Group
//   REset [CONFIG|STATS]
//         config    - reset configuration to defaults
//         stats     - reset cell statistics
//   ENable | DISable [BALANCE|CANTERM|C3100R|THermistor]
//         balance   - enable/disable cell balancing
//         canterm   - enable/disable CAN termination resistor
//         3100r     - enable/disable Curtis 3100R display
//         thermistor <ltc> <num> - enable/disable thermistor
//   LOCK            - lock configuration
//   UPGRADE         - performs a firmware upgrade
// bmsc> help

// bmsc> sh stats
// total|-mean cell voltage------|-standard deviation----------------
//      | 4.013v                 |  0.001v
// pack1|-voltage----min----max--|----deviation-----min---max--delta-
//  c1  | 4.011v   4.011v 4.013v | -0.002v  -1.6s  -1.8s -1.3s  0.5s
//  c2  | 4.013v   4.013v 4.014v |  0.000v  +0.0s  +0.0s +0.0s  0.0s
//  c3  | 4.013v   4.012v 4.014v | -0.000v  +0.0s  +0.0s +0.0s  0.0s
//  c4  | 4.013v   4.013v 4.014v |  0.000v  +0.0s  +0.0s +0.0s  0.0s
//  c5  | 4.014v   4.013v 4.015v |  0.001v  +0.0s  +0.0s +0.0s  0.0s
//  c6  | 4.014v   4.013v 4.015v |  0.001v  +1.1s  +0.0s +1.2s  1.2s
//  c7  | 4.013v   4.012v 4.014v | -0.000v  +0.0s  +0.0s +0.0s  0.0s
//  c8  | 4.013v   4.013v 4.015v |  0.000v  +0.0s  +0.0s +0.0s  0.0s
//  c9  | 4.013v   4.013v 4.014v |  0.000v  +0.0s  +0.0s +0.0s  0.0s
//  c10 | 4.014v   4.014v 4.016v |  0.002v  +1.6s  +1.4s +2.0s  0.6s
//  c11 | 4.014v   4.014v 4.015v |  0.001v  +1.4s  +1.2s +1.6s  0.4s
//  c12 | 4.014v   4.014v 4.015v |  0.001v  +1.3s  +0.0s +1.4s  1.4s
//  c13 | 4.011v   4.010v 4.012v | -0.002v  -2.5s  -2.5s -2.0s  0.5s
//  c14 | 4.012v   4.011v 4.013v | -0.001v  -1.5s  -1.8s -1.3s  0.5s
//  c15 | 4.012v   4.012v 4.013v | -0.001v  +0.0s  +0.0s +0.0s  0.0s
//  c16 | 4.013v   4.013v 4.015v |  0.001v  +0.0s  +0.0s +0.0s  0.0s
//  c17 | 4.013v   4.012v 4.014v | -0.000v  +0.0s  +0.0s +0.0s  0.0s
//  c18 | 4.012v   4.011v 4.013v | -0.001v  -1.2s  -1.3s +0.0s  1.3s
//  c19 | 4.013v   4.012v 4.014v | -0.000v  +0.0s  +0.0s +0.0s  0.0s
//  c20 | 4.013v   4.012v 4.014v | -0.000v  +0.0s  +0.0s +0.0s  0.0s
//  c21 | 4.014v   4.013v 4.015v |  0.001v  +0.0s  +0.0s +0.0s  0.0s
//  c22 | 4.014v   4.013v 4.015v |  0.001v  +1.1s  +0.0s +1.5s  1.5s
//  c23 | 4.014v   4.013v 4.015v |  0.001v  +1.1s  +0.0s +1.4s  1.4s
//  c24 | 4.013v   4.013v 4.015v |  0.001v  +0.0s  +0.0s +0.0s  0.0s
// bmsc> sh stats
type help struct {
}
