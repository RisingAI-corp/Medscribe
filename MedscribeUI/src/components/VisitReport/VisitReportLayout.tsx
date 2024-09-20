import { useState } from 'react';
import { Flex } from '@mantine/core';
import SoapSectionBox from './SoapSectionBox/SoapSectionBox';

function VisitReportLayout() {
  // TEMPORARY DATA
  const [soapData, setSoapData] = useState([
    {
      title: 'Summary',
      text: 'Patient reports that they are feeling well.',
      isLoading: true,
      timestamp: '1 min ago',
    },
    {
      title: 'Subjective',
      text: 'Patient reports that they are feeling well.',
      isLoading: true,
    },
    {
      title: 'Objective',
      text: 'Vitals are within normal range.',
      isLoading: false,
    },
    {
      title: 'Assessment',
      text: 'Patient is in good health.',
      isLoading: false,
    },
    {
      title: 'Plan',
      text: 'Continue with current medication and follow up in 6 months.',
      isLoading: false,
    },
  ]);

  const handleSoapDataUpdate = (index: number, newText: string) => {
    setSoapData(
      soapData.map((section, i) => {
        if (i === index) {
          return { ...section, text: newText };
        }
        return section;
      }),
    );
  };

  const handleLearnFormat = () => {
    // TO BE IMPLEMENTED
  };

  return (
    <Flex direction="column" gap="xl">
      {soapData.map((section, index) => (
        <SoapSectionBox
          key={index}
          title={section.title}
          text={section.text}
          isLoading={section.isLoading}
          timestamp={section.timestamp} // Pass timestamp if present
          handleSave={(newText: string) => {
            handleSoapDataUpdate(index, newText);
          }}
          handleLearnFormat={handleLearnFormat}
        />
      ))}
    </Flex>
  );
}

export default VisitReportLayout;
