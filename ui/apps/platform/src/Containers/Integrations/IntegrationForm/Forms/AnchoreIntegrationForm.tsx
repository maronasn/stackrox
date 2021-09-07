import React, { ReactElement } from 'react';
import { TextInput, PageSection, Form, Checkbox } from '@patternfly/react-core';
import * as yup from 'yup';

import usePageState from 'Containers/Integrations/hooks/usePageState';
import useIntegrationForm from '../useIntegrationForm';
import { IntegrationFormProps } from '../integrationFormTypes';

import IntegrationFormActions from '../IntegrationFormActions';
import FormCancelButton from '../FormCancelButton';
import FormTestButton from '../FormTestButton';
import FormSaveButton from '../FormSaveButton';
import FormMessage from '../FormMessage';
import FormLabelGroup from '../FormLabelGroup';

export type AnchoreIntegration = {
    id?: string;
    name: string;
    categories: 'REGISTRY'[];
    anchore: {
        endpoint: string;
        username: string;
        password: string;
        insecure: boolean;
    };
    skipTestIntegration: boolean;
    type: 'anchore';
    enabled: boolean;
    clusterIds: string[];
};

export type AnchoreIntegrationFormValues = {
    config: AnchoreIntegration;
    updatePassword: boolean;
};

export const validationSchema = yup.object().shape({
    config: yup.object().shape({
        name: yup.string().trim().required('An integration name is required'),
        categories: yup
            .array()
            .of(yup.string().trim().oneOf(['REGISTRY']))
            .min(1, 'Must have at least one type selected')
            .required('A category is required'),
        anchore: yup.object().shape({
            endpoint: yup.string().trim().required('An endpoint is required'),
            username: yup.string().trim(),
            password: yup
                .string()
                .test(
                    'password-test',
                    'A password is required',
                    (value, context: yup.TestContext) => {
                        const requirePasswordField =
                            // eslint-disable-next-line @typescript-eslint/ban-ts-comment
                            // @ts-ignore
                            context?.from[2]?.value?.updatePassword || false;

                        if (!requirePasswordField) {
                            return true;
                        }

                        const trimmedValue = value?.trim();
                        return !!trimmedValue;
                    }
                ),
            insecure: yup.bool(),
        }),
        skipTestIntegration: yup.bool(),
        type: yup.string().matches(/anchore/),
        enabled: yup.bool(),
        clusterIds: yup.array().of(yup.string()),
    }),
    updatePassword: yup.bool(),
});

export const defaultValues: AnchoreIntegrationFormValues = {
    config: {
        name: '',
        categories: ['REGISTRY'],
        anchore: {
            endpoint: '',
            username: '',
            password: '',
            insecure: false,
        },
        skipTestIntegration: false,
        type: 'anchore',
        enabled: true,
        clusterIds: [],
    },
    updatePassword: true,
};

function AnchoreIntegrationForm({
    initialValues = null,
    isEditable = false,
}: IntegrationFormProps<AnchoreIntegration>): ReactElement {
    const formInitialValues = defaultValues;
    if (initialValues) {
        formInitialValues.config = {
            ...formInitialValues.config,
            ...initialValues,
        };
        // We want to clear the password because backend returns '******' to represent that there
        // are currently stored credentials
        formInitialValues.config.anchore.password = '';
    }
    const {
        values,
        touched,
        errors,
        dirty,
        isValid,
        setFieldValue,
        handleBlur,
        isSubmitting,
        isTesting,
        onSave,
        onTest,
        onCancel,
        message,
    } = useIntegrationForm<AnchoreIntegrationFormValues, typeof validationSchema>({
        initialValues: formInitialValues,
        validationSchema,
    });
    const { isCreating } = usePageState();

    function onChange(value, event) {
        return setFieldValue(event.target.id, value);
    }

    return (
        <>
            <PageSection variant="light" isFilled hasOverflowScroll>
                {message && <FormMessage message={message} />}
                <Form isWidthLimited>
                    <FormLabelGroup
                        label="Integration name"
                        isRequired
                        fieldId="config.name"
                        touched={touched}
                        errors={errors}
                    >
                        <TextInput
                            isRequired
                            type="text"
                            id="config.name"
                            value={values.config.name}
                            onChange={onChange}
                            onBlur={handleBlur}
                            isDisabled={!isEditable}
                        />
                    </FormLabelGroup>
                    <FormLabelGroup
                        label="Endpoint"
                        isRequired
                        fieldId="config.anchore.endpoint"
                        touched={touched}
                        errors={errors}
                    >
                        <TextInput
                            isRequired
                            type="text"
                            id="config.anchore.endpoint"
                            value={values.config.anchore.endpoint}
                            onChange={onChange}
                            onBlur={handleBlur}
                            isDisabled={!isEditable}
                        />
                    </FormLabelGroup>
                    <FormLabelGroup
                        label="Username"
                        fieldId="config.anchore.username"
                        touched={touched}
                        errors={errors}
                    >
                        <TextInput
                            isRequired
                            type="text"
                            id="config.anchore.username"
                            value={values.config.anchore.username}
                            onChange={onChange}
                            onBlur={handleBlur}
                            isDisabled={!isEditable}
                        />
                    </FormLabelGroup>
                    {!isCreating && (
                        <FormLabelGroup
                            fieldId="updatePassword"
                            helperText="Setting this to false will use the currently stored credentials, if they exist."
                            errors={errors}
                        >
                            <Checkbox
                                label="Update stored credentials"
                                id="updatePassword"
                                aria-label="update password"
                                isChecked={values.updatePassword}
                                onChange={onChange}
                                onBlur={handleBlur}
                                isDisabled={!isEditable}
                            />
                        </FormLabelGroup>
                    )}
                    <FormLabelGroup
                        isRequired={values.updatePassword}
                        label="Password"
                        fieldId="config.anchore.password"
                        touched={touched}
                        errors={errors}
                    >
                        <TextInput
                            isRequired={values.updatePassword}
                            type="password"
                            id="config.anchore.password"
                            value={values.config.anchore.password}
                            onChange={onChange}
                            onBlur={handleBlur}
                            isDisabled={!isEditable || !values.updatePassword}
                            placeholder={
                                values.updatePassword
                                    ? ''
                                    : 'Currently-stored password will be used.'
                            }
                        />
                    </FormLabelGroup>
                    <FormLabelGroup
                        fieldId="config.anchore.insecure"
                        touched={touched}
                        errors={errors}
                    >
                        <Checkbox
                            label="Disable TLS certificate validation (insecure)"
                            id="config.anchore.insecure"
                            isChecked={Boolean(values.config.anchore.insecure)}
                            onChange={onChange}
                            onBlur={handleBlur}
                            isDisabled={!isEditable}
                        />
                    </FormLabelGroup>
                    <FormLabelGroup
                        fieldId="config.skipTestIntegration"
                        touched={touched}
                        errors={errors}
                    >
                        <Checkbox
                            label="Create integration without testing"
                            id="config.skipTestIntegration"
                            isChecked={Boolean(values.config.skipTestIntegration)}
                            onChange={onChange}
                            onBlur={handleBlur}
                            isDisabled={!isEditable}
                        />
                    </FormLabelGroup>
                </Form>
            </PageSection>
            {isEditable && (
                <IntegrationFormActions>
                    <FormSaveButton
                        onSave={onSave}
                        isSubmitting={isSubmitting}
                        isTesting={isTesting}
                        isDisabled={!dirty || !isValid}
                    >
                        Save
                    </FormSaveButton>
                    <FormTestButton
                        onTest={onTest}
                        isSubmitting={isSubmitting}
                        isTesting={isTesting}
                        isValid={isValid}
                    >
                        Test
                    </FormTestButton>
                    <FormCancelButton onCancel={onCancel}>Cancel</FormCancelButton>
                </IntegrationFormActions>
            )}
        </>
    );
}

export default AnchoreIntegrationForm;
