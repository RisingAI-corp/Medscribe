import { useState } from 'react';
import SoapSectionBox from './SoapSectionBox/SoapSectionBox';

function VisitReportLayout() {
  const [soapData, setSoapData] = useState([
    { title: 'Subjective', text: 'Patient reports that they are feeling well.' },
    { title: 'Objective', text: 'Vitals are within normal range.' },
    { title: 'Assessment', text: 'Patient is in good health.' },
    { title: 'Plan', text: 'Continue with current medication and follow up in 6 months.' },
  ]);

  return (
    <>
      {console.log('rerender')}
      {soapData.map((section, index) => (
        <SoapSectionBox
          key={index}
          title={section.title}
          text={section.text}
          onTextChange={newText => {
            const updatedSoapData = [...soapData];
            updatedSoapData[index] = { ...updatedSoapData[index], text: newText };
            setSoapData(updatedSoapData);
          }}
        />
      ))}
    </>
  );
}

export default VisitReportLayout;
