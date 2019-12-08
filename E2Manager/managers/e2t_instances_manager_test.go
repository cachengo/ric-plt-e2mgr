package managers

import (
	"e2mgr/configuration"
	"e2mgr/logger"
	"e2mgr/mocks"
	"e2mgr/services"
	"fmt"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

const E2TAddress = "10.10.2.15:9800"
const E2TAddress2 = "10.10.2.16:9800"

func initE2TInstancesManagerTest(t *testing.T) (*mocks.RnibReaderMock, *mocks.RnibWriterMock, *E2TInstancesManager) {
	logger, err := logger.InitLogger(logger.DebugLevel)
	if err != nil {
		t.Errorf("#... - failed to initialize logger, error: %s", err)
	}
	config := &configuration.Configuration{RnibRetryIntervalMs: 10, MaxRnibConnectionAttempts: 3}

	readerMock := &mocks.RnibReaderMock{}
	writerMock := &mocks.RnibWriterMock{}
	rnibDataService := services.NewRnibDataService(logger, config, readerMock, writerMock)
	e2tInstancesManager := NewE2TInstancesManager(rnibDataService, logger)
	return readerMock, writerMock, e2tInstancesManager
}

func TestAddNewE2TInstanceSaveE2TInstanceFailure(t *testing.T) {
	rnibReaderMock, rnibWriterMock, e2tInstancesManager := initE2TInstancesManagerTest(t)
	rnibWriterMock.On("SaveE2TInstance", mock.Anything).Return(common.NewInternalError(errors.New("Error")))
	err := e2tInstancesManager.AddE2TInstance(E2TAddress)
	assert.NotNil(t, err)
	rnibReaderMock.AssertNotCalled(t, "GetE2TAddresses")
}

func TestAddNewE2TInstanceGetE2TAddressesInternalFailure(t *testing.T) {
	rnibReaderMock, rnibWriterMock, e2tInstancesManager := initE2TInstancesManagerTest(t)
	rnibWriterMock.On("SaveE2TInstance", mock.Anything).Return(nil)
	e2tAddresses := []string{}
	rnibReaderMock.On("GetE2TAddresses").Return(e2tAddresses, common.NewInternalError(errors.New("Error")))
	err := e2tInstancesManager.AddE2TInstance(E2TAddress)
	assert.NotNil(t, err)
	rnibReaderMock.AssertNotCalled(t, "SaveE2TAddresses")
}

func TestAddNewE2TInstanceNoE2TAddresses(t *testing.T) {
	rnibReaderMock, rnibWriterMock, e2tInstancesManager := initE2TInstancesManagerTest(t)
	rnibWriterMock.On("SaveE2TInstance", mock.Anything).Return(nil)
	e2tAddresses := []string{}
	rnibReaderMock.On("GetE2TAddresses").Return(e2tAddresses, common.NewResourceNotFoundError(""))
	e2tAddresses = append(e2tAddresses, E2TAddress)
	rnibWriterMock.On("SaveE2TAddresses", e2tAddresses).Return(nil)
	err := e2tInstancesManager.AddE2TInstance(E2TAddress)
	assert.Nil(t, err)
	rnibWriterMock.AssertCalled(t, "SaveE2TAddresses", e2tAddresses)
}

func TestAddNewE2TInstanceEmptyE2TAddresses(t *testing.T) {
	rnibReaderMock, rnibWriterMock, e2tInstancesManager := initE2TInstancesManagerTest(t)
	rnibWriterMock.On("SaveE2TInstance", mock.Anything).Return(nil)
	e2tAddresses := []string{}
	rnibReaderMock.On("GetE2TAddresses").Return(e2tAddresses, nil)
	e2tAddresses = append(e2tAddresses, E2TAddress)
	rnibWriterMock.On("SaveE2TAddresses", e2tAddresses).Return(nil)
	err := e2tInstancesManager.AddE2TInstance(E2TAddress)
	assert.Nil(t, err)
	rnibWriterMock.AssertCalled(t, "SaveE2TAddresses", e2tAddresses)
}

func TestAddNewE2TInstanceSaveE2TAddressesFailure(t *testing.T) {
	rnibReaderMock, rnibWriterMock, e2tInstancesManager := initE2TInstancesManagerTest(t)
	rnibWriterMock.On("SaveE2TInstance", mock.Anything).Return(nil)
	E2TAddresses := []string{}
	rnibReaderMock.On("GetE2TAddresses").Return(E2TAddresses, nil)
	E2TAddresses = append(E2TAddresses, E2TAddress)
	rnibWriterMock.On("SaveE2TAddresses", E2TAddresses).Return(common.NewResourceNotFoundError(""))
	err := e2tInstancesManager.AddE2TInstance(E2TAddress)
	assert.NotNil(t, err)
}

func TestGetE2TInstanceSuccess(t *testing.T) {
	rnibReaderMock, _, e2tInstancesManager := initE2TInstancesManagerTest(t)
	address := "10.10.2.15:9800"
	e2tInstance := entities.NewE2TInstance(address)
	rnibReaderMock.On("GetE2TInstance", address).Return(e2tInstance, nil)
	res, err := e2tInstancesManager.GetE2TInstance(address)
	assert.Nil(t, err)
	assert.Equal(t, e2tInstance, res)
}

func TestGetE2TInstanceFailure(t *testing.T) {
	rnibReaderMock, _, e2tInstancesManager := initE2TInstancesManagerTest(t)
	var e2tInstance *entities.E2TInstance
	rnibReaderMock.On("GetE2TInstance", E2TAddress).Return(e2tInstance, common.NewInternalError(fmt.Errorf("for test")))
	res, err := e2tInstancesManager.GetE2TInstance(E2TAddress)
	assert.NotNil(t, err)
	assert.Nil(t, res)
}

func TestAssociateRanSuccess(t *testing.T) {
	rnibReaderMock, rnibWriterMock, e2tInstancesManager := initE2TInstancesManagerTest(t)
	e2tInstance  := entities.NewE2TInstance(E2TAddress)
	rnibReaderMock.On("GetE2TInstance", E2TAddress).Return(e2tInstance, nil)

	updateE2TInstance := *e2tInstance
	updateE2TInstance.AssociatedRanList = append(updateE2TInstance.AssociatedRanList, "test1")

	rnibWriterMock.On("SaveE2TInstance", &updateE2TInstance).Return(nil)

	err := e2tInstancesManager.AssociateRan("test1", E2TAddress)
	assert.Nil(t, err)
	rnibReaderMock.AssertExpectations(t)
	rnibWriterMock.AssertExpectations(t)
}

func TestAssociateRanGetInstanceFailure(t *testing.T) {
	rnibReaderMock, rnibWriterMock, e2tInstancesManager := initE2TInstancesManagerTest(t)

	var e2tInstance1 *entities.E2TInstance
	rnibReaderMock.On("GetE2TInstance", E2TAddress).Return(e2tInstance1, common.NewInternalError(fmt.Errorf("for test")))

	err := e2tInstancesManager.AssociateRan("test1", E2TAddress)
	assert.NotNil(t, err)
	rnibWriterMock.AssertNotCalled(t, "SaveE2TInstance")
}

func TestAssociateRanSaveInstanceFailure(t *testing.T) {
	rnibReaderMock, rnibWriterMock, e2tInstancesManager := initE2TInstancesManagerTest(t)
	e2tInstance1  := entities.NewE2TInstance(E2TAddress)
	rnibReaderMock.On("GetE2TInstance", E2TAddress).Return(e2tInstance1, nil)
	rnibWriterMock.On("SaveE2TInstance", mock.Anything).Return(common.NewInternalError(fmt.Errorf("for test")))

	err := e2tInstancesManager.AssociateRan("test1", E2TAddress)
	assert.NotNil(t, err)
	rnibReaderMock.AssertExpectations(t)
	rnibWriterMock.AssertExpectations(t)
}

func TestDissociateRanSuccess(t *testing.T) {
	rnibReaderMock, rnibWriterMock, e2tInstancesManager := initE2TInstancesManagerTest(t)

	e2tInstance := entities.NewE2TInstance(E2TAddress)
	e2tInstance.AssociatedRanList = []string{"test0","test1"}
	updatedE2TInstance := *e2tInstance
	updatedE2TInstance.AssociatedRanList = []string{"test0"}
	rnibReaderMock.On("GetE2TInstance", E2TAddress).Return(e2tInstance, nil)
	rnibWriterMock.On("SaveE2TInstance", &updatedE2TInstance).Return(nil)

	err := e2tInstancesManager.DissociateRan("test1", E2TAddress)
	assert.Nil(t, err)
	rnibReaderMock.AssertExpectations(t)
	rnibWriterMock.AssertExpectations(t)
}

func TestDissociateRanGetInstanceFailure(t *testing.T) {
	rnibReaderMock, rnibWriterMock, e2tInstancesManager := initE2TInstancesManagerTest(t)

	var e2tInstance1 *entities.E2TInstance
	rnibReaderMock.On("GetE2TInstance", E2TAddress).Return(e2tInstance1, common.NewInternalError(fmt.Errorf("for test")))
	err := e2tInstancesManager.DissociateRan("test1", E2TAddress)
	assert.NotNil(t, err)
	rnibWriterMock.AssertNotCalled(t, "SaveE2TInstance")
}

func TestDissociateRanSaveInstanceFailure(t *testing.T) {
	rnibReaderMock, rnibWriterMock, e2tInstancesManager := initE2TInstancesManagerTest(t)

	e2tInstance1  := entities.NewE2TInstance(E2TAddress)
	rnibReaderMock.On("GetE2TInstance", E2TAddress).Return(e2tInstance1, nil)
	rnibWriterMock.On("SaveE2TInstance", mock.Anything).Return(common.NewInternalError(fmt.Errorf("for test")))

	err := e2tInstancesManager.DissociateRan("test1", E2TAddress)
	assert.NotNil(t, err)
	rnibReaderMock.AssertExpectations(t)
	rnibWriterMock.AssertExpectations(t)
}

func TestSelectE2TInstancesGetE2TAddressesFailure(t *testing.T) {
	rnibReaderMock, rnibWriterMock, e2tInstancesManager := initE2TInstancesManagerTest(t)

	rnibReaderMock.On("GetE2TAddresses").Return([]string{}, common.NewInternalError(fmt.Errorf("for test")))
	address, err := e2tInstancesManager.SelectE2TInstance()
	assert.NotNil(t, err)
	assert.Empty(t, address)
	rnibReaderMock.AssertExpectations(t)
	rnibWriterMock.AssertNotCalled(t, "GetE2TInstances")
}

func TestSelectE2TInstancesEmptyE2TAddressList(t *testing.T) {
	rnibReaderMock, rnibWriterMock, e2tInstancesManager := initE2TInstancesManagerTest(t)

	rnibReaderMock.On("GetE2TAddresses").Return([]string{}, nil)
	address, err := e2tInstancesManager.SelectE2TInstance()
	assert.NotNil(t, err)
	assert.Empty(t, address)
	rnibReaderMock.AssertExpectations(t)
	rnibWriterMock.AssertNotCalled(t, "GetE2TInstances")
}

func TestSelectE2TInstancesGetE2TInstancesFailure(t *testing.T) {
	rnibReaderMock, rnibWriterMock, e2tInstancesManager := initE2TInstancesManagerTest(t)

	addresses := []string{E2TAddress}
	rnibReaderMock.On("GetE2TAddresses").Return(addresses, nil)
	rnibReaderMock.On("GetE2TInstances",addresses ).Return([]*entities.E2TInstance{}, common.NewInternalError(fmt.Errorf("for test")))
	address, err := e2tInstancesManager.SelectE2TInstance()
	assert.NotNil(t, err)
	assert.Empty(t, address)
	rnibReaderMock.AssertExpectations(t)
	rnibWriterMock.AssertExpectations(t)
}

func TestSelectE2TInstancesEmptyE2TInstancesList(t *testing.T) {
	rnibReaderMock, rnibWriterMock, e2tInstancesManager := initE2TInstancesManagerTest(t)

	addresses := []string{E2TAddress}
	rnibReaderMock.On("GetE2TAddresses").Return(addresses, nil)
	rnibReaderMock.On("GetE2TInstances",addresses ).Return([]*entities.E2TInstance{}, nil)
	address, err := e2tInstancesManager.SelectE2TInstance()
	assert.NotNil(t, err)
	assert.Empty(t, address)
	rnibReaderMock.AssertExpectations(t)
	rnibWriterMock.AssertExpectations(t)
}

func TestSelectE2TInstancesNoActiveE2TInstance(t *testing.T) {
	rnibReaderMock, rnibWriterMock, e2tInstancesManager := initE2TInstancesManagerTest(t)
	addresses := []string{E2TAddress,E2TAddress2}
	e2tInstance1 := entities.NewE2TInstance(E2TAddress)
	e2tInstance1.State = entities.ToBeDeleted
	e2tInstance1.AssociatedRanList = []string{"test1","test2","test3"}
	e2tInstance2 := entities.NewE2TInstance(E2TAddress2)
	e2tInstance2.State = entities.ToBeDeleted
	e2tInstance2.AssociatedRanList = []string{"test4","test5","test6", "test7"}

	rnibReaderMock.On("GetE2TAddresses").Return(addresses, nil)
	rnibReaderMock.On("GetE2TInstances",addresses).Return([]*entities.E2TInstance{e2tInstance1, e2tInstance2}, nil)
	address, err := e2tInstancesManager.SelectE2TInstance()
	assert.NotNil(t, err)
	assert.Equal(t, "", address)
	rnibReaderMock.AssertExpectations(t)
	rnibWriterMock.AssertExpectations(t)
}

func TestSelectE2TInstancesSuccess(t *testing.T) {
	rnibReaderMock, rnibWriterMock, e2tInstancesManager := initE2TInstancesManagerTest(t)
	addresses := []string{E2TAddress,E2TAddress2}
	e2tInstance1 := entities.NewE2TInstance(E2TAddress)
	e2tInstance1.AssociatedRanList = []string{"test1","test2","test3"}
	e2tInstance2 := entities.NewE2TInstance(E2TAddress2)
	e2tInstance2.AssociatedRanList = []string{"test4","test5","test6", "test7"}

	rnibReaderMock.On("GetE2TAddresses").Return(addresses, nil)
	rnibReaderMock.On("GetE2TInstances",addresses).Return([]*entities.E2TInstance{e2tInstance1, e2tInstance2}, nil)
	address, err := e2tInstancesManager.SelectE2TInstance()
	assert.Nil(t, err)
	assert.Equal(t, E2TAddress, address)
	rnibReaderMock.AssertExpectations(t)
	rnibWriterMock.AssertExpectations(t)
}

func TestRemoveE2TInstance(t *testing.T) {
	_, _, e2tInstancesManager := initE2TInstancesManagerTest(t)
	e2tInstance1  := entities.NewE2TInstance(E2TAddress)
	err := e2tInstancesManager.RemoveE2TInstance(e2tInstance1)
	assert.Nil(t, err)
}