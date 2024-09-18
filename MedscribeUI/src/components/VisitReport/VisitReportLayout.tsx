import { useState } from 'react';
import { Flex } from '@mantine/core';
import SoapSectionBox from './SoapSectionBox/SoapSectionBox';

function VisitReportLayout() {
  // TEMPORARY DATA
  const [soapData, setSoapData] = useState([
    { title: 'Summary', text: 'Patient reports that they are feeling well.', isSaving: true, timestamp: '1 min ago' },
    { title: 'Subjective', text: 'Patient reports that they are feeling well.', isSaving: true },
    { title: 'Objective', text: 'Vitals are within normal range.', isSaving: false },
    { title: 'Assessment', text: 'Patient is in good health.', isSaving: false },
    { title: 'Plan', text: 'Continue with current medication and follow up in 6 months.', isSaving: false },
  ]);

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
          isSaving={section.isSaving}
          timestamp={section.timestamp} // Pass timestamp if present
          handleSave={newText => {
            setSoapData(soapData.map((section, i) => (i === index ? { ...section, text: newText } : section)));
          }}
          handleLearnFormat={handleLearnFormat}
        />
      ))}
    </Flex>
  );
}

export default VisitReportLayout;
