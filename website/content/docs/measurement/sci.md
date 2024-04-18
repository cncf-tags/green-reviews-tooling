---
title: "Software Carbon Intensity"
description: "This section provides details for calculation and measurement of Software Carbon Intensity (SCI) specification."
summary: ""
date: 2023-09-07T16:04:48+02:00
lastmod: 2023-09-07T16:04:48+02:00
draft: false
slug: sci
weight: 810
toc: true
seo:
  title: "" # custom title (optional)
  description: "" # custom description (recommended)
  canonical: "" # custom canonical URL (optional)
  noindex: false # false (default) or true
---

The [Software Carbon Intensity (SCI) specification][1] by the [Green Software Foundation](https://greensoftware.foundation) is a method for **scoring** a software system based on the components that emit carbon.

The SCI can be used to reduce the total carbon footprint of software. However it is a _rate_ of carbon emissions for software and not necessarily the total carbon footprint.

The top-level formula of the SCI is the following:

> SCI = ([Energy (E)](#energy-e) * [Emissions Factor (I)](#energy-carbon-intensity-i)) + [Embodied Emissions (M)] per [Functional Unit (R)][5]

## Energy (E)

Energy (E) represents the operational energy consumed by the physical hardware that the software system operates on. More info in the [SCI Guide on (E)][2].

> Unit: kilowatt hours (kWh)

### Runtime energy

This project surfaces runtime energy consumption through [Kepler ( (Kubernetes-based Efficient Power Level Exporter)](https://github.com/sustainable-computing-io/kepler).

Kepler deploys a Kepler Exporter which runs as a Daemonset in the cluster. The Kepler Exporter scrapes energy metrics and exports them to Prometheus.

The metric chosen to calculate the energy consumption of Falco is [`kepler_container_joules_total (Counter)`](https://sustainable-computing.io/design/metrics):

> This metric is the aggregated package/socket energy consumption of CPU, dram, gpus, and other host components for a given container. Each component has individual metrics which are detailed next in this document.

In the [SCI dashboard](https://github.com/cncf-tags/green-reviews-tooling/blob/3a6266ceae99f40aaa367174ffb899385caf1d50/clusters/base/falco-sci.yaml#L505), the joules total is then converted to kWh using a hidden [dashboard variable](https://grafana.com/docs/grafana/latest/dashboards/variables) `$watt_per_second_to_kWh`. This standard unit conversion metric is equal to `0.000000277777777777778`.

## Energy Carbon Intensity (I)

The carbon intensity of electricity is a measure of how much carbon (CO2eq) emissions are produced per kilowatt-hour (kWh) of electricity consumed. More info in the [SCI Guide on (I)][3].

- The green reviews cluster is physically located in Equinix Metal's [Paris metro](https://deploy.equinix.com/locations).
- For carbon intensity we use the 2023 annual average value for France from the [CO2.js](https://github.com/thegreenwebfoundation/co2.js/blob/main/data/output/average-intensities.json#L422-L427) library from The Green Web Foundation.
- CO2.js uses data published by [Ember](https://ember-climate.org/data) under a CC-BY-4.0 license.

> Unit: carbon per kilowatt hours (gCO2eq/kWh)

### Power Usage Effectiveness (PUE)

Is a ratio that describes how efficiently a data center uses energy for computing. Specifically, how much energy is used by the computing equipment in contrast to cooling and any other overhead that supports the data center.

PUE = Total Facility Energy / IT Equipment Energy =  1 + (Non IT Facility Energy / IT Equipment Energy)

An ideal PUE is 1.0. Anything that isn't considered a computing device in a data center (e.g. lighting, cooling, etc.) falls into the category of facility energy consumption.

Equinix published this information for their data centers as part of their latest [2022 sustainability report](https://sustainability.equinix.com/wp-content/uploads/2023/05/Equinix-Inc.-2022-Sustainability-Report-Highlights-1.pdf).

## Embodied (M)

Embodied carbon (also known as embedded carbon) is the amount of carbon emitted during the creation and disposal of a hardware device. More info in the [SCI Guide on (M)][4].

> Unit: grams of carbon (gCO2eq)

```M = TE * (TiR/EL) * RS```

- TE = Total Embodied Emissions; the sum of Life Cycle Assessment (LCA) emissions for all hardware components.
- TiR = Time Reserved; the length of time the hardware is reserved for use by the software.
- EL = Expected Lifespan; the anticipated time that the equipment will be installed.
- RS = Resource-share; the share of the total available resources of the hardware reserved for use by the software.

- TE we use the [Boavizta API](https://doc.api.boavizta.org) passing the server spec published by Equinix Metal
- EL we use 4 years (35,040 hours) which is the value recommended by GSF
- TiR is 15 minutes (the duration of the green review)
- RS is 1 (we use bare metal servers so 100% of the resources are allocated to the software)

The Boavizta API uses data published by the manufacturers. You can read about their methodology in this [blog post](https://www.boavizta.org/en/blog/empreinte-de-la-fabrication-d-un-serveur).

### M Per Instance Type

| Instance Type	| TE (kgCO2eq) | TiR (minutes) | EL (years) | M (gCO2eq) |
|---------------|--------------|---------------|------------|------------|
| m3.small.x86  | 550          | 15            | 4          |	3.92     |

**m3.small.x86**

| Component   | Configuration        | Notes                                           |
|-------------|----------------------|-------------------------------------------------|
| CPU         | 1 Intel Xeon E-2378G | 8 cores @ 2.8 GHz, TDP of 80W, Rocket Lake arch |
| RAM         | 64 GB                |                                                 |
| SSD         | 2 x 480GB            |                                                 |
| Server Type | Rack                 |                                                 |
| PSU Quantity| 2                    |                                                 |

Server Spec: Equinix Metal [m3.small.x86](https://deploy.equinix.com/product/servers/m3-small)
CPU: Intel Xeon E-2378G [datasheet](https://www.intel.com/content/www/us/en/products/sku/212262/intel-xeon-e2378g-processor-16m-cache-2-80-ghz/specifications.html)

## Functional Unit (R)

The SCI is the rate of carbon emissions per one functional unit. The functional unit describes how the software application scales e.g. per additional user, API call, etc. Since each software application scales differently, each one also has a different functional unit. More information can be found in the [SCI Guide on (R)][5].

The SCI specification includes benchmark tests as a [suggested functional unit](https://sci.greensoftware.foundation/#functional-unit). The benchmark tests generate kernel events that Falco reacts to. The functional unit of the benchmark tests will be to reach a target kernel event rate e.g. X kernel events for Y minute(s).

<!-- Sources -->
[1]: https://sci.greensoftware.foundation
[2]: https://sci-guide.greensoftware.foundation/E
[3]: https://sci-guide.greensoftware.foundation/I
[4]: https://sci-guide.greensoftware.foundation/M
[5]: https://sci-guide.greensoftware.foundation/R
