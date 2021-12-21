import * as React from "react";
import App from "./table";
// styles
const pageStyles = {
  color: "#232129",
  padding: 96,
  fontFamily: "-apple-system, Roboto, sans-serif, serif",
};
const footerStyles = {
  // backgroundImage: `url('data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHZlcnNpb249IjEuMSIgdmlld0JveD0iMCAwIDc4NC4zNyAxMjc3LjM5Ij4KIDxnPgogICA8Zz4KICAgIDxwb2x5Z29uIGZpbGw9IiMzNDM0MzQiIGZpbGwtcnVsZT0ibm9uemVybyIgcG9pbnRzPSIzOTIuMDcsMCAzODMuNSwyOS4xMSAzODMuNSw4NzMuNzQgMzkyLjA3LDg4Mi4yOSA3ODQuMTMsNjUwLjU0ICIvPgogICAgPHBvbHlnb24gZmlsbD0iIzhDOEM4QyIgZmlsbC1ydWxlPSJub256ZXJvIiBwb2ludHM9IjM5Mi4wNywwIC0wLDY1MC41NCAzOTIuMDcsODgyLjI5IDM5Mi4wNyw0NzIuMzMgIi8+CiAgICA8cG9seWdvbiBmaWxsPSIjM0MzQzNCIiBmaWxsLXJ1bGU9Im5vbnplcm8iIHBvaW50cz0iMzkyLjA3LDk1Ni41MiAzODcuMjQsOTYyLjQxIDM4Ny4yNCwxMjYzLjI4IDM5Mi4wNywxMjc3LjM4IDc4NC4zNyw3MjQuODkgIi8+CiAgICA8cG9seWdvbiBmaWxsPSIjOEM4QzhDIiBmaWxsLXJ1bGU9Im5vbnplcm8iIHBvaW50cz0iMzkyLjA3LDEyNzcuMzggMzkyLjA3LDk1Ni41MiAtMCw3MjQuODkgIi8+CiAgICA8cG9seWdvbiBmaWxsPSIjMTQxNDE0IiBmaWxsLXJ1bGU9Im5vbnplcm8iIHBvaW50cz0iMzkyLjA3LDg4Mi4yOSA3ODQuMTMsNjUwLjU0IDM5Mi4wNyw0NzIuMzMgIi8+CiAgICA8cG9seWdvbiBmaWxsPSIjMzkzOTM5IiBmaWxsLXJ1bGU9Im5vbnplcm8iIHBvaW50cz0iMCw2NTAuNTQgMzkyLjA3LDg4Mi4yOSAzOTIuMDcsNDcyLjMzICIvPgogICA8L2c+CiAgPC9nPgo8L3N2Zz4K')`,
  bottom: 0,
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
}
const headingStyles = {
  marginTop: 0,
  marginBottom: 20,
  maxWidth: 1000,
};
const headingAccentStyles = {
  color: "#663399",
};
const paragraphStyles = {
  marginBottom: 48,
};
const codeStyles = {
  color: "#8A6534",
  padding: 4,
  backgroundColor: "#FFF4DB",
  fontSize: "1.25rem",
  borderRadius: 4,
};
const listStyles = {
  marginBottom: 96,
  paddingLeft: 0,
};
const listItemStyles = {
  fontWeight: 300,
  fontSize: 24,
  maxWidth: 560,
  marginBottom: 30,
};

const linkStyle = {
  color: "#8954A8",
  fontWeight: "bold",
  fontSize: 16,
  verticalAlign: "5%",
};

const docLinkStyle = {
  ...linkStyle,
  listStyleType: "none",
  marginBottom: 24,
};

const descriptionStyle = {
  color: "#232129",
  fontSize: 14,
  marginTop: 10,
  marginBottom: 0,
  lineHeight: 1.25,
};

const docLink = {
  text: "Chainflow",
  url: "https://chainflow.io/",
  color: "#8954A8",
};

// const badgeStyle = {
//   color: "#fff",
//   backgroundColor: "#088413",
//   border: "1px solid #088413",
//   fontSize: 11,
//   fontWeight: "bold",
//   letterSpacing: 1,
//   borderRadius: 4,
//   padding: "4px 6px",
//   display: "inline-block",
//   position: "relative",
//   top: -2,
//   marginLeft: 10,
//   lineHeight: 1,
// }

// data
const links = [
  {
    text: "Tutorial",
    url: "https://www.gatsbyjs.com/docs/tutorial/",
    description:
      "A great place to get started if you're new to web development. Designed to guide you through setting up your first Gatsby site.",
    color: "#E95890",
  },
  {
    text: "How to Guides",
    url: "https://www.gatsbyjs.com/docs/how-to/",
    description:
      "Practical step-by-step guides to help you achieve a specific goal. Most useful when you're trying to get something done.",
    color: "#1099A8",
  },
  {
    text: "Reference Guides",
    url: "https://www.gatsbyjs.com/docs/reference/",
    description:
      "Nitty-gritty technical descriptions of how Gatsby works. Most useful when you need detailed information about Gatsby's APIs.",
    color: "#BC027F",
  },
  {
    text: "Conceptual Guides",
    url: "https://www.gatsbyjs.com/docs/conceptual/",
    description:
      "Big-picture explanations of higher-level Gatsby concepts. Most useful for building understanding of a particular topic.",
    color: "#0D96F2",
  },
  {
    text: "Plugin Library",
    url: "https://www.gatsbyjs.com/plugins",
    description:
      "Add functionality and customize your Gatsby site or app with thousands of plugins built by our amazing developer community.",
    color: "#8EB814",
  },
  {
    text: "Build and Host",
    url: "https://www.gatsbyjs.com/cloud",
    badge: true,
    description:
      "Now you’re ready to show the world! Give your Gatsby site superpowers: Build and host on Gatsby Cloud. Get started for free!",
    color: "#663399",
  },
];

// markup
const IndexPage = () => {
  return (
    <div>
      <main style={pageStyles}>
        <title>Nakamoto Coefficients - Chainflow</title>
        <h1 style={headingStyles}>
          <span style={headingAccentStyles}> Nakamoto Coefficient </span>
        </h1>
        <ul style={listStyles}>
          {/* <li style={docLinkStyle}>
          <a
            style={linkStyle}
            href={`${docLink.url}?utm_source=starter&utm_medium=start-page&utm_campaign=minimal-starter`}
          >
            {docLink.text}
          </a>
        </li> */}
          <App />
          {/* {links.map((link) => (
          <li key={link.url} style={{ ...listItemStyles, color: link.color }}>
            <span>
              <a
                style={linkStyle}
                href={`${link.url}?utm_source=starter&utm_medium=start-page&utm_campaign=minimal-starter`}
              >
                {link.text}
              </a>
              {link.badge && (
                <span style={badgeStyle} aria-label="New Badge">
                  NEW!
                </span>
              )}
              <p style={descriptionStyle}>{link.description}</p>
            </span>
          </li>
        ))} */}
        </ul>
      </main>
      <footer style={footerStyles}>
        <p>
          Developed with &hearts; by{" "}
          <a href="https://twitter.com/xenowits">xenowits</a>{" / "}Sponsored by{" "}
          <a href="https://chainflow.io/">Chainflow</a>
        </p>
      </footer>
    </div>
  );
};

export default IndexPage;
