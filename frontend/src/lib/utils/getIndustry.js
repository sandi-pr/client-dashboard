// Define the data for each industry
export const industryData = [
  {
    id: 1,
    industry: "Tech Industry",
    points: [
      {
        title: "Automated Compliance",
        desc: "Documents AI ensures compliance by monitoring regulatory documents and policies, flagging any changes or risks."
      },
      {
        title: "Product Documentation",
        desc: "Analyze extensive product documentation to generate summaries, update logs, and highlight key product specs."
      },
      {
        title: "Client Proposal",
        desc: "Streamline the creation of client proposals by extracting and reusing relevant content from previous documents in the Document Library."
      }
    ]
  },
  {
    id: 2,
    industry: "Banking Industry",
    points: [
      {
        title: "Sop Compliance",
        desc: "Documents AI reviews SOPs and other procedural documents, checking adherence to internal and external regulations."
      },
      {
        title: "Contract Agreement Summarization",
        desc: "Summarize complex financial contracts and highlight critical clauses for faster review."
      },
      {
        title: "Customer Data Privacy",
        desc: "Analyze data protection policies and documentation to ensure compliance with privacy regulations like GDPR."
      }
    ]
  },
  {
    id: 3,
    industry: "Oil & Gas Industry",
    points: [
      {
        title: "Safety Protocols",
        desc: "Easily access and analyze safety procedure documents, ensuring compliance with safety standards across locations."
      },
      {
        title: "Equipment Maintenance",
        desc: "Store, access, and analyze equipment maintenance logs, helping predict maintenance schedules and reduce downtime."
      },
      {
        title: "Environmental Compliance",
        desc: "Track and analyze documents on environmental regulations to ensure adherence and minimize environmental impact."
      }
    ]
  },
  {
    id: 4,
    industry: "Retail Industry",
    points: [
      {
        title: "Inventory Management",
        desc: "Track inventory management SOPs and policies to maintain smooth logistics and restocking processes."
      },
      {
        title: "Sales Performance",
        desc: "Analyze sales reports, identify trends, and generate insights to improve promotional and pricing strategies."
      },
      {
        title: "Customer Experience SOPs",
        desc: "Review and optimize customer experience guidelines based on customer service documentation and feedback."
      }
    ]
  },
  {
    id: 5,
    industry: "Healthcare Industry",
    points: [
      {
        title: "Patient Care",
        desc: "Centralize patient care protocols, ensuring staff compliance and improving treatment consistency."
      },
      {
        title: "Medical Research",
        desc: "Summarize complex medical research documents to help healthcare professionals stay updated on the latest findings."
      },
      {
        title: "Compliance Documentation",
        desc: "Ensure compliance with healthcare standards by monitoring policy and procedural documents for updates."
      }
    ]
  },
  {
    id: 6,
    industry: "Education Industry",
    points: [
      {
        title: "Curriculum Planning",
        desc: "Review and summarize curriculum documents, helping educators align courses with academic standards."
      },
      {
        title: "Student Performance",
        desc: "Track and analyze student performance reports and feedback to identify areas needing improvement."
      },
      {
        title: "Policy Compliance",
        desc: "Ensure institutional compliance with education policies by monitoring and updating policy documents."
      }
    ]
  },
  {
    id: 7,
    industry: "Financial Services Industry",
    points: [
      {
        title: "Risk Assessment",
        desc: "Analyze and summarize risk assessment documents, providing quick insights into potential risk areas."
      },
      {
        title: "Policy Regulation",
        desc: "Documents AI monitors evolving regulations and policy changes, ensuring compliance across services."
      },
      {
        title: "Investment Research",
        desc: "Generate summaries of financial research documents, giving insights into investment opportunities."
      }
    ]
  },
  {
    id: 8,
    industry: "Hospitality Industry",
    points: [
      {
        title: "Service Standards",
        desc: "Store and update service standards to ensure a consistent guest experience across all locations."
      },
      {
        title: "Employee Training",
        desc: "Optimize training programs by storing and updating training materials for easy employee access."
      },
      {
        title: "Customer Feedback",
        desc: "Process and analyze guest feedback to improve service quality and identify key customer preferences."
      }
    ]
  },
  {
    id: 9,
    industry: "Manufacturing Industry",
    points: [
      {
        title: "Quality Control",
        desc: "Store and track quality control procedures to ensure products meet regulatory and internal standards."
      },
      {
        title: "Production Report",
        desc: "Summarize production reports to identify bottlenecks, resource needs, and efficiency improvements."
      },
      {
        title: "Equipment Maintenance",
        desc: "Monitor maintenance logs and predict upcoming maintenance needs to reduce production downtime."
      }
    ]
  },
  {
    id: 10,
    industry: "Logistics Industry",
    points: [
      {
        title: "Supply Chain",
        desc: "Centralize supply chain documents, enabling faster analysis and optimization of logistics processes."
      },
      {
        title: "Route Delivery",
        desc: "Analyze delivery and route reports to improve efficiency and reduce transportation costs."
      },
      {
        title: "Regulatory Compliance",
        desc: "Ensure compliance with transportation and logistics regulations by tracking related documentation."
      }
    ]
  },
  {
    id: 11,
    industry: "Telecommunications Industry",
    points: [
      {
        title: "Network Maintenance",
        desc: "Track and review network maintenance logs to predict maintenance needs and avoid service disruptions."
      },
      {
        title: "Customer Service Script",
        desc: "Analyze customer service documents to improve scripts and enhance customer support quality."
      },
      {
        title: "Regulatory Compliance",
        desc: "Monitor compliance with telecom regulations by reviewing and updating related documentation."
      }
    ]
  },
  {
    id: 12,
    industry: "Real Estate Industry",
    points: [
      {
        title: "Lease Agreement",
        desc: "Summarize lengthy lease agreements, highlighting essential terms for quick reference."
      },
      {
        title: "Property Management",
        desc: "Centralize maintenance logs and property records for easy access, improving efficiency in property management."
      },
      {
        title: "Market Analysis",
        desc: "Analyze real estate market reports, providing insights to guide investment and marketing strategies."
      }
    ]
  }
];

// Function to retrieve points by industry name
export function getIndustryPoints(industryName) {
  // Find the industry data matching the industryName
  const selectedIndustry = industryData.find(industry => industry.industry === industryName);
  
  // Check if industry is found, then return points, else return a message
  if (selectedIndustry) {
    return selectedIndustry.points.map(point => ({
      title: point.title,
      desc: point.desc
    }));
  } else {
    console.log(`Industry "${industryName}" not found.`);
    return;
  }
}

// Example usage
// const selectedIndustryPoints = getIndustryPoints("Tech Industry");
// console.log(selectedIndustryPoints);
