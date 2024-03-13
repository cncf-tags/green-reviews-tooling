# Software Carbon Intensity

SCI = ([E](https://sci-guide.greensoftware.foundation/E) * [I](https://sci-guide.greensoftware.foundation/I) * [PUE](https://learn.greensoftware.foundation/energy-efficiency/#power-usage-effectiveness)) + [M](https://sci-guide.greensoftware.foundation/M) per [R](https://sci-guide.greensoftware.foundation/R)

## Energy (E)

TBC

## Energy Carbon Intensity (I)

TBC

## Power Usage Efectiveness (PUE)
Is a ratio that describes how efficiently a data center uses energy for computing. Specifically, how much energy is used by the computing equipment in contrast to cooling and any other overhead that supports the data center.

PUE = Total Facility Energy / IT Equipment Energy =  1 + Non IT Facility Energy / IT Equipment Energy

An ideal PUE is 1.0. Anything that isn't considered a computing device in a data center (e.g. lighting, cooling, etc.) falls into the category of facility energy consumption.

Equinix published this information for their data centers as part of their [last sustainability report](https://sustainability.equinix.com/wp-content/uploads/2023/05/Equinix-Inc.-2022-Sustainability-Report-Highlights-1.pdf).

## Embodied (M)

Embodied carbon (also known as embedded carbon) is the amount of carbon emitted during the creation and disposal of a hardware device.

M = TE * (TiR/EL) * RS

- TE = Total Embodied Emissions; the sum of Life Cycle Assessment (LCA) emissions for all hardware components.
- TiR = Time Reserved; the length of time the hardware is reserved for use by the software.
- EL = Expected Lifespan; the anticipated time that the equipment will be installed.
- RS = Resource-share; the share of the total available resources of the hardware reserved for use by the software.

More info [here](https://sci-guide.greensoftware.foundation/M).

- TE we use the [Boavizta API](https://doc.api.boavizta.org/) passing the server spec published by Equinix Metal 
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

Server Spec: Equinix Metal [m3.small.x86](https://deploy.equinix.com/product/servers/m3-small/)
CPU: Intel Xeon E-2378G [datasheet](https://www.intel.com/content/www/us/en/products/sku/212262/intel-xeon-e2378g-processor-16m-cache-2-80-ghz/specifications.html)

## Functional Unit (R)

TBC
