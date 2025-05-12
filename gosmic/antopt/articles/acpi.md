---
title: "Does ACPI dream of electric sheep, during S5?"
date: "2025-03-21"
description: "What does State After G3 mean? Do I need to set it to S5 or S0 State? This article is a brief overview of ACPI global and sleep states." 
categories:
  - "computers"
published: true
---
I bought a mini PC and I wanted to set it up so that if power goes out, it automatically starts again when power is back. In the BIOS settings I had to change this configuration:

**State After G3**: could be set to **S5 State** or **S0 State**.

But what does this mean exactly?

## ACPI

Before we get into G3, S5 or S0, let's start with an introduction of who is responsible for deciding what they mean. 

The "Advanced Configuration and Power Interface" is the open standard that can be found in every computer nowadays. Initially developed in December 1996 by Intel, Microsoft, and Toshiba, and now owned by the [UEFI Forum](https://en.wikipedia.org/wiki/UEFI_Forum) alliance.
As you can imagine, having a bunch of hardware manufacturers and software developers to agree on a common interface is not a joke. ACPI is one of those things that I never really understood or cared about, but simplified my life indirectly.

## Global states (Gx)

Once your OS activates ACPI, it's their responsibility to take care of everything related to power management.

ACPI defines four global states. Your ACPI-compliant machine is always in one of these states:

- **G0** (Working): your CPU is running and executing instructions.
- **G1** (Sleeping): it depends on the sleep state.
- **G2** (Soft Off): your machine has been powered down.
- **G3** (Mechanical Off): your machine is not receiving any power (note: your machine still probably has a small battery that keeps the internal real-time clock running).

## Sleep states

Sleeping mode comes with different flavors, so ACPI defines six different states. Your ACPI-compliant machine is always in one of these sleep states:

- **S0**: not sleeping. If your machine is on G0, it's also on S0.
- **S0i1**, **S0i2**, **S0i3**: Intel and ARM substates of S0, where some power consumption is reduced but going back to the normal S0 only takes a couple milliseconds so it's unnoticeable for the user.
- **S1**: CPU stops executing instructions and its caches are flushed. RAM is maintained. Other peripherals (e.g. monitors) might turn off.
- **S2**: CPU is completely powered off.
- **S3**: the common "standby". RAM remains in a low-power mode and pretty much everything else is powered off.
- **S4**: referred to as "hibernation". The RAM is copied over to the hard drive before shutting down the system (i.e. entering S5).
- **S5**: shutdown. If you machine is on G2, it's also on S5.

It should be clearer now what that configuration in the BIOS settings meant by "State after G3".
It's a way of setting what state you want to set your machine to when passing from G3 to G2 (i.e. when connecting a source of power). In my case, I set it to S0, which means to start working.

Ciao!

P.S. we still often refer to "BIOS settings", even it should be technically more accurate to call it "UEFI settings" or more generically "firmware settings".
