/*
 * CR3 Client Dashboard
 * Package: gitlab.hexacode.org/cr3/client-dashboard
 * Maintainer: Azzis Arswendo <azzis@hexacode.org>
 * 
 * Copyright (C) 2024 Hexacode Teknologi Indonesia
 * All Rights Reserved
 */

/** @type {import('tailwindcss').Config}*/
const config = {
  content: [
    "./src/**/*.{html,js,svelte,ts}",
    "./node_modules/flowbite-svelte/**/*.{html,js,svelte,ts}",
  ],

  plugins: [require("flowbite/plugin"), require("@tailwindcss/typography"), require("tailwind-fluid-typography")],

  darkMode: "class",

  theme: {
    fluidTypography: {},
    fontFamily: {
      jakarta: ["Plus Jakarta Sans", "sans-serif"],
    },
    extend: {
      colors: {
        primary: {
          25: "#FCFAFF",
          50: "#F9F5FF",
          100: "#F4EBFF",
          200: "#E9D7FE",
          300: "#D6BBFB",
          400: "#B692F6",
          500: "#9E77ED",
          600: "#7F56D9",
          700: "#6941C6",
          800: "#53389E",
          900: "#42307D",
          950: "#2C1C5F",
        },
        hc: {
          graylight: {
            25: "#FCFCFD",
            50: "#F9FAFB",
            100: "#F2F4F7",
            200: "#EAECF0",
            300: "#D0D5DD",
            400: "#98A2B3",
            500: "#667085",
            600: "#475467",
            700: "#344054",
            800: "#182230",
            900: "#101828",
            950: "#0C111D",
          },
          graydark: {
            25: "#FAFAFA",
            50: "#F5F5F6",
            100: "#F0F1F1",
            200: "#ECECED",
            300: "#CECFD2",
            400: "#94969C",
            500: "#85888E",
            600: "#61646C",
            700: "#333741",
            800: "#1F242F",
            900: "#161B26",
            950: "#0C111D",
          },
          brand: {
            25: "#FCFAFF",
            50: "#F9F5FF",
            100: "#F4EBFF",
            200: "#E9D7FE",
            300: "#D6BBFB",
            400: "#B692F6",
            500: "#9E77ED",
            600: "#7F56D9",
            700: "#6941C6",
            800: "#53389E",
            900: "#42307D",
            950: "#2C1C5F",
          },
          error: {
            25: "#FFFBFA",
            50: "#FEF3F2",
            100: "#FEE4E2",
            200: "#FECDCA",
            300: "#FDA29B",
            400: "#F97066",
            500: "#F04438",
            600: "#D92D20",
            700: "#B42318",
            800: "#912018",
            900: "#7A271A",
            950: "#55160C",
          },
          warning: {
            25: "#FFFCF5",
            50: "#FFFAEB",
            100: "#FEF0C7",
            200: "#FEDF89",
            300: "#FEC84B",
            400: "#FDB022",
            500: "#F79009",
            600: "#DC6803",
            700: "#B54708",
            800: "#93370D",
            900: "#7A2E0E",
            950: "#4E1D09",
          },
          success: {
            25: "#F6FEF9",
            50: "#ECFDF3",
            100: "#DCFAE6",
            200: "#ABEFC6",
            300: "#75E0A7",
            400: "#47CD89",
            500: "#17B26A",
            600: "#079455",
            700: "#067647",
            800: "#085D3A",
            900: "#074D31",
            950: "#053321",
          },
        },
        fg: {
          primary: {
            light: "#101828",
            dark: "#FFFFFF",
          },
          secondary: {
            light: "#344054",
            dark: "#CECFD2",
            hover: {
              light: "#182230",
              dark: "#ECECED",
            },
          },
          tertiary: {
            light: "#475467",
            dark: "#94969C",
            hover: {
              light: "#344054",
              dark: "#CECFD2",
            },
          },
          quaternary: {
            light: "#667085",
            dark: "#94969C",
            hover: {
              light: "#475467",
              dark: "#CECFD2",
            },
          },
          quinary: {
            light: "#98A2B3",
            dark: "#85888E",
            hover: {
              light: "#667085",
              dark: "#94969C",
            },
          },
          senary: {
            light: "#D0D5DD",
            dark: "#61646C",
          },
          disabled: {
            light: "#98A2B3",
            dark: "#85888E",
            subtle: {
              light: "#D0D5DD",
              dark: "#61646C",
            },
          },
          brand: {
            primary: {
              light: "#7F56D9",
              dark: "#9E77ED",
              alt: {
                light: "#7F56D9",
                dark: "#CECFD2",
              },
            },
            secondary: {
              light: "#9E77ED",
              dark: "#9E77ED",
            },
          },
          error: {
            primary: {
              light: "#D92D20",
              dark: "#F04438",
            },
            secondary: {
              light: "#F04438",
              dark: "#F97066",
            },
          },
          warning: {
            primary: {
              light: "#DC6803",
              dark: "#F79009",
            },
            secondary: {
              light: "#F79009",
              dark: "#FDB022",
            },
          },
          success: {
            primary: {
              light: "#079455",
              dark: "#17B26A",
            },
            secondary: {
              light: "#17B26A",
              dark: "#47CD89",
            },
          },
        },
        utility: {
          gray: {
            50: {
              light: "#F9FAFB",
              dark: "#161B26",
            },
            100: {
              light: "#F2F4F7",
              dark: "#1F242F",
            },
            200: {
              light: "#EAECF0",
              dark: "#333741",
            },
            300: {
              light: "#D0D5DD",
              dark: "#61646C",
            },
            400: {
              light: "#98A2B3",
              dark: "#475467",
            },
            500: {
              light: "#667085",
              dark: "#667085",
            },
            600: {
              light: "#475467",
              dark: "#98A2B3",
            },
            700: {
              light: "#344054",
              dark: "#D0D5DD",
            },
            800: {
              light: "#182230",
              dark: "#ECECED",
            },
            900: {
              light: "#101828",
              dark: "#F0F1F1",
            },
          },
          brand: {
            50: {
              light: "#F9F5FF",
              dark: "#2C1C5F",
            },
            100: {
              light: "#F4EBFF",
              dark: "#42307D",
            },
            200: {
              light: "#E9D7FE",
              dark: "#53389E",
            },
            300: {
              light: "#D6BBFB",
              dark: "#6941C6",
            },
            400: {
              light: "#B692F6",
              dark: "#7F56D9",
            },
            500: {
              light: "#9E77ED",
              dark: "#9E77ED",
            },
            600: {
              light: "#7F56D9",
              dark: "#B692F6",
            },
            700: {
              light: "#6941C6",
              dark: "#D6BBFB",
            },
            800: {
              light: "#53389E",
              dark: "#E9D7FE",
            },
            900: {
              light: "#42307D",
              dark: "#F4EBFF",
            },
          },
          error: {
            50: {
              light: "#FEF3F2",
              dark: "#7A271A",
            },
            100: {
              light: "#FEE4E2",
              dark: "#912018",
            },
            200: {
              light: "#FECDCA",
              dark: "#B42318",
            },
            300: {
              light: "#FDA29B",
              dark: "#D92D20",
            },
            400: {
              light: "#F97066",
              dark: "#F04438",
            },
            500: {
              light: "#F04438",
              dark: "#F97066",
            },
            600: {
              light: "#D92D20",
              dark: "#FDA29B",
            },
            700: {
              light: "#B42318",
              dark: "#FECDCA",
            },
            800: {
              light: "#912018",
              dark: "#FEE4E2",
            },
            900: {
              light: "#7A271A",
              dark: "#FEF3F2",
            },
          },
          warning: {
            50: {
              light: "#FFFAEB",
              dark: "#7A2E0E",
            },
            100: {
              light: "#FEF0C7",
              dark: "#93370D",
            },
            200: {
              light: "#FEDF89",
              dark: "#B54708",
            },
            300: {
              light: "#FEC84B",
              dark: "#DC6803",
            },
            400: {
              light: "#FDB022",
              dark: "#F79009",
            },
            500: {
              light: "#F79009",
              dark: "#FDB022",
            },
            600: {
              light: "#DC6803",
              dark: "#FEC84B",
            },
            700: {
              light: "#B54708",
              dark: "#FEDF89",
            },
            800: {
              light: "#93370D",
              dark: "#FEF0C7",
            },
            900: {
              light: "#7A2E0E",
              dark: "#FFFAEB",
            },
          },
          success: {
            50: {
              light: "#ECFDF3",
              dark: "#074D31",
            },
            100: {
              light: "#DCFAE6",
              dark: "#085D3A",
            },
            200: {
              light: "#ABEFC6",
              dark: "#067647",
            },
            300: {
              light: "#75E0A7",
              dark: "#079455",
            },
            400: {
              light: "#47CD89",
              dark: "#17B26A",
            },
            500: {
              light: "#17B26A",
              dark: "#47CD89",
            },
            600: {
              light: "#079455",
              dark: "#75E0A7",
            },
            700: {
              light: "#067647",
              dark: "#ABEFC6",
            },
            800: {
              light: "#085D3A",
              dark: "#DCFAE6",
            },
            900: {
              light: "#074D31",
              dark: "#ECFDF3",
            },
          },
        },
      },
      textColor: {
        primary: {
          light: "#101828",
          dark: "#F5F5F6",
          onbrand: {
            light: "#FFFFFF",
            dark: "#F5F5F6",
          },
        },
        secondary: {
          light: "#344054",
          dark: "#CECFD2",
          hover: {
            light: "#182230",
            dark: "#CECFD2",
          },
          onbrand: {
            light: "#E9D7FE",
            dark: "#CECFD2",
          },
        },
        tertiary: {
          light: "#475467",
          dark: "#94969C",
          hover: {
            light: "#344054",
            dark: "#CECFD2",
          },
          onbrand: {
            light: "#E9D7FE",
            dark: "#94969C",
          },
        },
        quarterary: {
          light: "#667085",
          dark: "#94969C",
          onbrand: {
            light: "#E9D7FE",
            dark: "#94969C",
          },
        },
        disabled: {
          light: "#667085",
          dark: "#85888E",
        },
        placeholder: {
          light: "#667085",
          dark: "#94969C",
          subtle: {
            light: "#D0D5DD",
            dark: "#333741",
          },
        },
        brand: {
          primary: {
            light: "#42307D",
            dark: "#F5F5F6",
          },
          secondary: {
            light: "#6941C6",
            dark: "#CECFD2",
          },
          tertiary: {
            light: "#7F56D9",
            dark: "#94969C",
            alt: {
              light: "#7F56D9",
              dark: "#F5F5F6",
            },
          },
        },
        error: {
          primary: {
            light: "#D92D20",
            dark: "#F97066",
          },
        },
        warning: {
          primary: {
            light: "#DC6803",
            dark: "#FDB022",
          },
        },
        success: {
          primary: {
            light: "#079455",
            dark: "#FDB022",
          },
        },
      },
      borderColor: {
        primary: {
          light: "#D0D5DD",
          dark: "#333741",
        },
        secondary: {
          light: "#EAECF0",
          dark: "#1F242F",
        },
        tertiary: {
          light: "#F2F4F7",
          dark: "#1F242F",
        },
        disabled: {
          light: "#D0D5DD",
          dark: "#333741",
          subtle: {
            light: "#EAECF0",
            dark: "#1F242F",
          },
        },
        brand: {
          light: "#D6BBFB",
          dark: "#B692F6",
          solid: {
            light: "#7F56D9",
            dark: "#9E77ED",
            alt: {
              light: "#7F56D9",
              dark: "#333741",
            },
          },
        },
        error: {
          light: "#FDA29B",
          dark: "#F97066",
          solid: {
            light: "#D92D20",
            dark: "#F04438",
          },
        },
      },
      backgroundColor: {
        primary: {
          light: "#FFFFFF",
          dark: "#0C111D",
          alt: {
            light: "#FFFFFF",
            dark: "#161B26",
          },
          hover: {
            light: "#F9FAFB",
            dark: "#1F242F",
          },
          solid: {
            light: "#0C111D",
            dark: "#161B26",
          },
        },
        secondary: {
          light: "#F9FAFB",
          dark: "#161B26",
          alt: {
            light: "#F9FAFB",
            dark: "#0C111D",
          },
          hover: {
            light: "#F2F4F7",
            dark: "#1F242F",
          },
          subtle: {
            light: "#FCFCFD",
            dark: "#161B26",
          },
          solid: {
            light: "#475467",
            dark: "#61646C",
          },
        },
        tertiary: {
          light: "#F2F4F7",
          dark: "#1F242F",
        },
        quarterary: {
          light: "#EAECF0",
          dark: "#333741",
        },
        active: {
          light: "#F9FAFB",
          dark: "#1F242F",
        },
        disabled: {
          light: "#F2F4F7",
          dark: "#1F242F",
          subtle: {
            light: "#F9FAFB",
            dark: "#161B26",
          },
        },
        overlay: {
          light: "#0C111D",
          dark: "#1F242F",
        },
        brand: {
          primary: {
            light: "#F9F5FF",
            dark: "#9E77ED",
            alt: {
              light: "#F9F5FF",
              dark: "#1F242F",
            },
          },
          secondary: {
            light: "#F4EBFF",
            dark: "#7F56D9",
          },
          solid: {
            light: "#7F56D9",
            dark: "#7F56D9",
            hover: {
              light: "#6941C6",
              dark: "#9E77ED",
            },
          },
          section: {
            light: "#53389E",
            dark: "#1F242F",
            subtle: {
              light: "#6941C6",
              dark: "#0C111D",
            },
          },
        },
        error: {
          primary: {
            light: "#FEF3F2",
            dark: "#F04438",
          },
          secondary: {
            light: "#FEE4E2",
            dark: "#D92D20",
          },
          solid: {
            light: "#D92D20",
            dark: "#D92D20",
          },
          warning: {
            primary: {
              light: "#FFFAEB",
              dark: "#F79009",
            },
            secondary: {
              light: "#FEF0C7",
              dark: "#DC6803",
            },
            solid: {
              light: "#DC6803",
              dark: "#DC6803",
            },
          },
          success: {
            primary: {
              light: "#ECFDF3",
              dark: "#17B26A",
            },
            secondary: {
              light: "#DCFAE6",
              dark: "#079455",
            },
            solid: {
              light: "#079455",
              dark: "#079455",
            },
          },
        },
      },
    },
  },
};

module.exports = config;
