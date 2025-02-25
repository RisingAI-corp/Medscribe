import { Button } from '@mantine/core';
import PatientRegistryLayout from '../../components/PatientRegistry/patientRegistryLayout';
import PatientDashBoardLayout from '../../components/PaitentDashboard/patientDashBoardLayout';
import Header from '../../components/Header/header';
import { currentlySelectedPatientAtom } from '../../states/patientsAtom';
import PatientReception from '../../components/PateintReception/patientReceptionLayout';
import { useAtom } from 'jotai';

function HomeScreen() {
  const [patientId, setPatientId] = useAtom(currentlySelectedPatientAtom);
  return (
    <div className="flex flex-col flex-1">
      <Header />

      <div className="flex flex-1 flex-row overflow-hidden">
        <div className="flex flex-col w-[400px] border-r border-gray-300 h-full box-border">
          <div className="p-4 border-b border-gray-300 bg-white">
            <Button
              variant="outline"
              className="border-blue-500 text-blue-500 font-bold w-full"
              onClick={() => {
                setPatientId('');
              }}
            >
              + START A VISIT
            </Button>
          </div>
          <div className="flex-1 p-4 overflow-y-auto">
            <PatientRegistryLayout />
          </div>
        </div>
        {patientId !== '' ? (
          <div className="flex-1 h-full overflow-y-auto box-border">
            <PatientDashBoardLayout key={patientId} />
          </div>
        ) : (
          <div className="flex-1 h-screen flex justify-center items-center">
            <PatientReception />
          </div>
        )}
      </div>
    </div>
  );
}

export default HomeScreen;
