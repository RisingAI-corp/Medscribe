import React from 'react';

interface PatientBackgroundDetailsProps {
  name: string;
  condensedSummary: string;
  lastVisitDate: string;
  duration: number;
  lastVisitSummary: string;
}

const PatientBackgroundDetails: React.FC<PatientBackgroundDetailsProps> = ({
  name,
  condensedSummary,
  lastVisitDate,
  duration,
  lastVisitSummary,
}) => {
  // Format duration from milliseconds to minutes
  const durationInMinutes = Math.round(duration / 60000);
  
  return (
    <div className="p-6 bg-transparent">
      <div className="mb-4">
        <h2 className="text-2xl font-bold text-gray-800">Patient Background</h2>
        <p className="text-sm text-gray-500">{name}</p>
      </div>

      <div className="grid grid-cols-2 gap-4 mb-4">
        <div>
          <h3 className="text-lg font-semibold text-gray-700">Last Visit</h3>
          <p className="mt-1 text-gray-600">{new Date(lastVisitDate).toLocaleDateString()}</p>
        </div>
        <div>
          <h3 className="text-lg font-semibold text-gray-700">Duration</h3>
          <p className="mt-1 text-gray-600">{durationInMinutes} minutes</p>
        </div>
      </div>
      
      <div className="mb-4">
        <h3 className="text-lg font-semibold text-gray-700">Condensed Summary</h3>
        <p className="mt-1 text-gray-600">{condensedSummary}</p>
      </div>
      
      
      
      <div>
        <h3 className="text-lg font-semibold text-gray-700">Last Visit Summary</h3>
        <p className="mt-1 text-gray-600">{lastVisitSummary}</p>
      </div>
    </div>
  );
};

export default PatientBackgroundDetails; 