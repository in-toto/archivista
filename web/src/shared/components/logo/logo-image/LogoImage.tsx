/* eslint-disable @typescript-eslint/no-unsafe-argument */
import { GatsbyImage, getImage } from 'gatsby-plugin-image';
import { graphql, useStaticQuery } from 'gatsby';

import React from 'react';

const Logo = () => {
  const { logoFile } = useStaticQuery(graphql`
    query Logo {
      logoFile: file(relativePath: { glob: "icon-transparent.png" }) {
        childImageSharp {
          gatsbyImageData(width: 32, placeholder: BLURRED, formats: [AUTO, WEBP, AVIF])
        }
      }
    }
  `);

  const logoImage = getImage(logoFile);

  return <>{logoImage && <GatsbyImage image={logoImage} alt="Judge logo" style={{ height: 32 }} />}</>;
};

export default Logo;
